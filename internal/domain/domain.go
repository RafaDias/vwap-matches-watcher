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
