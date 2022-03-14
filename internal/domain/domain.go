// Package domain is the hearth for project
// It provides a principal layer.
package domain

import (
	"errors"
	"time"
)

type Price struct {
	Price float64
	Size  float64
}

type Transaction struct {
	ProductID string    `json:"product_id"`
	Price     float64   `json:"price"`
	Size      string    `json:"size"`
	Time      time.Time `json:"time"`
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


func (tp *TradingPair) dropOldest() {
	oldestMatch := tp.Prices[0]
	tp.Size -= oldestMatch.Size
	tp.Volume -= oldestMatch.Size * oldestMatch.Price
}

func (tp *TradingPair) isFull() bool {
	return len(tp.Prices) == tp.WindowSize
}
