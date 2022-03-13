package coinbase

import (
	"github.com/rafadias/crypto-watcher/internal/application/providers/exchange"
	"github.com/rafadias/crypto-watcher/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

const filePath = "/Users/rafaeldias/projects/github.com/RafaDias/crypto-watcher/config/config.json"

func TestService_Connect_With_Correct_Values(t *testing.T) {
	cfg, err := config.FromPath(filePath)

	if err != nil {
		t.Fatal(err)
	}

	svc, err := New(exchange.Config{
		BaseUrl:       cfg.Exchange.BaseUrl,
		Channels:      cfg.Exchange.Channels,
		Subscriptions: cfg.Exchange.Subscriptions,
	})

	assert.Nil(t, err)
	assert.Equal(t, cfg.Exchange.Subscriptions, svc.GetSubscriptions())
}

func TestService_WrongConnectionMustReturnAnError(t *testing.T) {
	cfg, err := config.FromPath(filePath)

	if err != nil {
		t.Fatal(err)
	}

	_, err = New(exchange.Config{
		BaseUrl:       "wss:/wrong-host/with-wrong-path",
		Channels:      cfg.Exchange.Channels,
		Subscriptions: cfg.Exchange.Subscriptions,
	})

	assert.NotNil(t, err)
}
