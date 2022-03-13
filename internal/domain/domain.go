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

func (tp *TradingPair) isFull() bool {
	return len(tp.Prices) == tp.WindowSize
}