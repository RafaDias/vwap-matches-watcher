// Package exchange provides an exchange API.
package exchange

import (
	"github.com/rafadias/crypto-watcher/internal/domain"
)

type Service interface {
	ListenTransactions(chan domain.Transaction) error
	GetSubscriptions() []string
	GetChannels() []string
	Close() error
}

type Config struct {
	BaseURL       string
	Channels      []string
	Subscriptions []string
}
