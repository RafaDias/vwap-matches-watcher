package unit

import (
	"github.com/joho/godotenv"
	"github.com/rafadias/crypto-watcher/internal/application/providers/exchange"
	"github.com/rafadias/crypto-watcher/internal/config"
	"github.com/rafadias/crypto-watcher/internal/infra/providers/exchange/coinbase"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)


func TestService_Connect_With_Correct_Values(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Print("No .env file found")
	}

	cfg := config.New()

	svc, err := coinbase.New(exchange.Config{
		BaseUrl:       cfg.Exchange.BaseUrl,
		Channels:      cfg.Exchange.Channels,
		Subscriptions: cfg.Exchange.Subscriptions,
	})

	assert.Nil(t, err)
	assert.Equal(t, cfg.Exchange.Subscriptions, svc.GetSubscriptions())
}

func TestService_WrongConnectionMustReturnAnError(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	cfg := config.New()

	_, err := coinbase.New(exchange.Config{
		BaseUrl:       "wss:/wrong-host/with-wrong-path",
		Channels:      cfg.Exchange.Channels,
		Subscriptions: cfg.Exchange.Subscriptions,
	})

	assert.NotNil(t, err)
}
