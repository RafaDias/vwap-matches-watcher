package domain

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type Price struct {
	Price float64
	Size  float64
}

type TradingPair struct {
	Name       string
	Prices     []Price
	Size       float64
	Volume     float64
	WindowSize int
}

func (tp *TradingPair) VWAP() float64 {
	return tp.Volume / tp.Size
}

func (tp *TradingPair) Add(p Price) error {
	if p.Price <= 0 || p.Size <= 0 {
		return errors.New("you must provide positive number")
	}

	if tp.isFull() {
		tp.dropOldest()
	}
	tp.Prices = append(tp.Prices, p)
	tp.Size += p.Size
	tp.Volume += p.Size * p.Price

	return nil
}

func (tp *TradingPair) Listen(c chan Price, logger *log.Logger) {
	for msg := range c {
		err := tp.Add(msg)
		if err != nil {
		}
		logger.Println(fmt.Sprintf("Wrap %s: %f", tp.Name, tp.VWAP()))
	}
}

func (tp *TradingPair) dropOldest() {
	oldestMatch := tp.Prices[0]
	tp.Size -= oldestMatch.Size
	tp.Volume -= oldestMatch.Size * oldestMatch.Price
}

func (tp *TradingPair) isFull() bool {
	return len(tp.Prices) == tp.WindowSize
}