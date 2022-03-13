package unit

import (
	"github.com/rafadias/crypto-watcher/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTradingParShouldCalculateVWAPFromOneTradingPair(t *testing.T) {
	tp := domain.TradingPair{WindowSize: 1}
	err := tp.Add(domain.Price{Price: 1.0, Size: 2.0})
	assert.Equal(t, 1.0,  tp.VWAP())
	assert.Nil(t, err)
}


func TestTradingParShouldCalculateVWAPFromTwoTradingPairs(t *testing.T) {
	tp := domain.TradingPair{WindowSize: 2}
	err := tp.Add(domain.Price{Price: 4.0, Size: 2.0}) // (8 + 48) / 8
	assert.Nil(t, err)
	err = tp.Add(domain.Price{Price: 8.0, Size: 6.0})
	assert.Nil(t, err)
	assert.Equal(t,7.0,  tp.VWAP())
}

func TestTradingPairShouldRemoveTheOldestValue(t *testing.T) {
	tp := domain.TradingPair{WindowSize: 2}
	err := tp.Add(domain.Price{Price: 4.0, Size: 2.0}) // will be ignored, because window length is 2
	assert.Nil(t, err)
	err = tp.Add(domain.Price{Price: 8.0, Size: 6.0}) // (48 + 50) / 11
	assert.Nil(t, err)
	err = tp.Add(domain.Price{Price: 10.0, Size: 5.0})
	assert.Nil(t, err)
	assert.Equal(t, 8.909090909090908,  tp.VWAP())
}

func TestShouldNotAllowNegativeNumberOnPrice(t *testing.T) {
	tp := domain.TradingPair{WindowSize: 2}
	err := tp.Add(domain.Price{Price: -4.0, Size: 2.0})
	assert.NotNil(t, err)
}

func TestShouldNotAllowNegativeNumberOnSize(t *testing.T) {
	tp := domain.TradingPair{WindowSize: 2}
	err := tp.Add(domain.Price{Price: 4.0, Size: -2.0})
	assert.NotNil(t, err)
}
