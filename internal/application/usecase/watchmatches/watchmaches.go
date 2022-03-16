// Package watchmatches provides an usecase API.
package watchmatches

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/rafadias/crypto-watcher/internal/application/providers/exchange"
	"github.com/rafadias/crypto-watcher/internal/domain"
)

type watchMatcherUseCase struct {
	log       *log.Logger
	exchange  exchange.Service
	consumers map[string]consumer
	wg        *sync.WaitGroup
}

func New(log *log.Logger, exchange exchange.Service) *watchMatcherUseCase {
	svc := &watchMatcherUseCase{
		exchange: exchange,
		log:      log,
		wg:       new(sync.WaitGroup),
	}
	return svc
}

func (wm *watchMatcherUseCase) Execute(windowSize int, ctx context.Context) {
	wm.setupSubscriptions(windowSize)
	wm.watch(ctx)
}

func (wm *watchMatcherUseCase) Wait() {
	wm.wg.Wait()
}

func (wm *watchMatcherUseCase) GetVWAP() map[string]float64 {
	vwap := make(map[string]float64)
	for _, subscription := range wm.consumers {
		vwap[subscription.tradingPair.Name] = subscription.tradingPair.VWAP()
	}
	return vwap
}

func (wm *watchMatcherUseCase) watch(ctx context.Context) {
	transaction := make(chan domain.Transaction)

	go func() {
		err := wm.exchange.ListenTransactions(transaction)
		if err != nil {
			wm.log.Fatal(fmt.Sprintf("err %s: details %s", domain.ErrCannotListenExchange.Error(), err.Error()))
		}
	}()

	for {
		select {
		case txn, ok := <-transaction:
			if !ok {
				wm.log.Println("INFO: exchange channel is closed")
				wm.close()
				return
			}
			c, ok := wm.consumers[txn.ProductID]
			if !ok {
				wm.log.Fatal(domain.ErrCannotListenExchange.Error())
			}

			c.priceChannel <- domain.Price{Price: txn.Price, Size: txn.Price}

		case <-ctx.Done():
			wm.log.Println("WARN: Received cancellation signal, closing consumers!")
			wm.close()
			return
		}
	}
}

func (wm *watchMatcherUseCase) close() {
	for _, sub := range wm.consumers {
		close(sub.priceChannel)
	}
	wm.log.Println("INFO: consumers are closed")
}

func (wm *watchMatcherUseCase) setupSubscriptions(windowSize int) {
	tps := make(map[string]consumer)
	wm.wg.Add(len(wm.exchange.GetSubscriptions()))

	for _, name := range wm.exchange.GetSubscriptions() {
		tp := domain.TradingPair{
			Name:       name,
			WindowSize: windowSize,
		}
		incomingMatches := make(chan domain.Price)
		cm := consumer{tradingPair: &tp, priceChannel: incomingMatches}
		go func() {
			defer wm.wg.Done()
			defer wm.log.Printf("Finish goroutine: %d", consumerID())
			cm.listen(incomingMatches, wm.log)
		}()
		tps[name] = cm
	}
	wm.consumers = tps
}

func consumerID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
