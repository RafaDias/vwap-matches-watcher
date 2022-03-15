// Package watchmatches provides an usecase API.
package watchmatches

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/rafadias/crypto-watcher/internal/application/providers/exchange"
	"github.com/rafadias/crypto-watcher/internal/domain"
)

type subscription struct {
	priceChannel chan domain.Price
	tradingPair  *domain.TradingPair
}

func (s *subscription) listen(priceChan <-chan domain.Price, logger *log.Logger) {
	for msg := range priceChan {
		err := s.tradingPair.Add(msg)
		if err != nil {
			log.Printf("err %v trying to add item, it should compromise vwap value", err)
		}
		logger.Println(fmt.Sprintf("Wrap %s: %f", s.tradingPair.Name, s.tradingPair.VWAP()))
	}
}

type watchMatcherUseCase struct {
	log           *log.Logger
	exchange      exchange.Service
	subscriptions map[string]subscription
	wg            sync.WaitGroup
}

func New(log *log.Logger, exchange exchange.Service) *watchMatcherUseCase {
	svc := &watchMatcherUseCase{
		exchange: exchange,
		log:      log,
	}
	return svc
}

func (wm *watchMatcherUseCase) Execute(windowSize int) {
	wm.setupSubscriptions(windowSize)
	wm.watch()
}

func (wm *watchMatcherUseCase) watch() {
	transaction := make(chan domain.Transaction)
	wm.wg.Add(1)

	go func() {
		defer wm.wg.Done()
		err := wm.exchange.ListenTransactions(transaction)
		if err != nil {
			wm.log.Fatal("error accours", err)
		}
	}()

	for txn := range transaction {
		c, ok := wm.subscriptions[txn.ProductID]
		if !ok {
			wm.log.Fatal("not a valid channel")
		}
		size, err := strconv.ParseFloat(txn.Size, 64)
		if err != nil {
			wm.log.Fatal("error parsing value", err)
		}
		c.priceChannel <- domain.Price{Price: txn.Price, Size: size}
		wm.log.Println(txn)
	}

	for _, sub := range wm.subscriptions {
		close(sub.priceChannel)
	}
	wm.wg.Wait()
}

func (wm *watchMatcherUseCase) setupSubscriptions(windowSize int) {
	tps := make(map[string]subscription)
	wm.wg.Add(len(wm.exchange.GetSubscriptions()))

	for _, name := range wm.exchange.GetSubscriptions() {
		tp := domain.TradingPair{
			Name:       name,
			WindowSize: windowSize,
		}
		incomingMatches := make(chan domain.Price)
		sub := subscription{tradingPair: &tp, priceChannel: incomingMatches}
		go func() {
			defer wm.wg.Done()
			defer wm.log.Printf("Terminei a goroutinen do: %s", name)
			sub.listen(incomingMatches, wm.log)
		}()
		tps[name] = sub
	}
	wm.subscriptions = tps
}

func (wm *watchMatcherUseCase) GetVWAP() map[string]float64 {
	vwap := make(map[string]float64)
	for _, subscription := range wm.subscriptions {
		vwap[subscription.tradingPair.Name] = subscription.tradingPair.VWAP()
	}
	return vwap
}
