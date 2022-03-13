package memory

import (
	"github.com/rafadias/crypto-watcher/internal/application/providers/exchange"
	"github.com/rafadias/crypto-watcher/internal/domain"
)


type inmemoryExchange struct {
	subscriptions []string
	channels      []string
}

func (i *inmemoryExchange) GetChannels() []string {
	return i.subscriptions
}

func (i *inmemoryExchange) Close() error {
	return nil
}

func (i *inmemoryExchange) ListenTransactions(transactions chan domain.Transaction) error {
	txns := []domain.Transaction{
		{
			ProductId: i.subscriptions[0],
			Price:     10.00,
			Size:      "0.10",
		},
		{
			ProductId: i.subscriptions[1],
			Price:     20.00,
			Size:      "0.3",
		},
		{
			ProductId: i.subscriptions[2],
			Price:     10.00,
			Size:      "0.01",
		},
	}
	for _, t := range txns {
		transactions <- t
	}

	close(transactions)
	return nil
}

func (i *inmemoryExchange) GetSubscriptions() []string {
	return i.subscriptions
}

func New(subscriptions, channels []string) exchange.Service {
	return &inmemoryExchange{
		subscriptions: subscriptions,
		channels:      channels,
	}
}
