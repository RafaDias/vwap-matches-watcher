package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTradingParShouldCalculateVWAPFromOneTradingPair(t *testing.T) {
	tp := TradingPair{WindowSize: 1}
	err := tp.Add(Price{Price: 1.0, Size: 2.0})
	assert.Equal(t, 1.0,  tp.VWAP())
	assert.Nil(t, err)
}

