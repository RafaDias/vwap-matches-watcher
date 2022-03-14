// Package watchmatches provides an usecase API.
package watchmatches

import (
	"fmt"
	"log"
	"strconv"

	"github.com/rafadias/crypto-watcher/internal/application/providers/exchange"
	"github.com/rafadias/crypto-watcher/internal/domain"
)

type subscription struct {
	price       chan domain.Price
	tradingPair *domain.TradingPair
}

func (s *subscription) Listen(c chan domain.Price, logger *log.Logger) {
	for msg := range c {
		err := s.tradingPair.Add(msg)
		if err != nil {
			logger.Fatal("cannot add item")
		}
		logger.Println(fmt.Sprintf("Wrap %s: %f", s.tradingPair.Name, s.tradingPair.VWAP()))
	}
}

type watchMatcherUseCase struct {
	log           *log.Logger
	exchange      exchange.Service
	subscriptions map[string]subscription
}

func New(log *log.Logger, exchange exchange.Service, windowSize int) *watchMatcherUseCase {
	svc := &watchMatcherUseCase{
		exchange: exchange,
		log:      log,
	}
	svc.setupSubscriptions(windowSize)
	return svc
}

func (wm *watchMatcherUseCase) Execute() {
	wm.watch()
}

func (wm *watchMatcherUseCase) watch() {
	transaction := make(chan domain.Transaction)
	go func() {
		err := wm.exchange.ListenTransactions(transaction)
		if err != nil {
			wm.log.Fatal("error accours", err)
		}
	}()

	for txn := range transaction {
		c, ok := wm.subscriptions[txn.ProductID]
		if !ok {
			wm.log.Fatal("error")
		}
		size, err := strconv.ParseFloat(txn.Size, 64)
		if err != nil {
			wm.log.Fatal("error accours", err)
		}
		c.price <- domain.Price{Price: txn.Price, Size: size}
		wm.log.Println(txn)
	}
}

func (wm *watchMatcherUseCase) setupSubscriptions(windowSize int) {
	tps := make(map[string]subscription)
	for _, name := range wm.exchange.GetSubscriptions() {
		tp := domain.TradingPair{
			Name:       name,
			WindowSize: windowSize,
		}
		incomingMatches := make(chan domain.Price)
		sub := subscription{tradingPair: &tp, price: incomingMatches}
		go sub.Listen(incomingMatches, wm.log)
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
