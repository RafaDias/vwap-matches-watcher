package watchmatches

import (
	"fmt"
	"github.com/rafadias/crypto-watcher/internal/domain"
	"log"
)

type consumer struct {
	priceChannel chan domain.Price
	tradingPair  *domain.TradingPair
}

func (s *consumer) listen(priceChan <-chan domain.Price, logger *log.Logger) {
	for msg := range priceChan {
		err := s.tradingPair.Add(msg)
		if err != nil {
			log.Printf("err %v trying to add item, it should compromise vwap value", err)
		}
		logger.Println(fmt.Sprintf("VWAP %s: %f", s.tradingPair.Name, s.tradingPair.VWAP()))
	}
}
