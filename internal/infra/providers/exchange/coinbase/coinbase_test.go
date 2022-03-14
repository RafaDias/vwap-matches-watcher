package coinbase

import (
	"github.com/joho/godotenv"
	"github.com/rafadias/crypto-watcher/internal/application/providers/exchange"
	"github.com/rafadias/crypto-watcher/internal/config"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

const filePath = "../../../../../.env"

func TestService_Connect_With_Correct_Values(t *testing.T) {
	if err := godotenv.Load(filePath); err != nil {
		log.Print("No .env file found")
	}

	cfg := config.New()

	svc, err := New(exchange.Config{
		BaseURL:       cfg.Exchange.BaseURL,
		Channels:      cfg.Exchange.Channels,
		Subscriptions: cfg.Exchange.Subscriptions,
	})

	assert.Nil(t, err)
	assert.Equal(t, cfg.Exchange.Subscriptions, svc.GetSubscriptions())
}

func TestService_WrongConnectionMustReturnAnError(t *testing.T) {
	if err := godotenv.Load(filePath); err != nil {
		log.Print("No .env file found")
	}

	cfg := config.New()

	_, err := New(exchange.Config{
		BaseURL:       "wss:/wrong-host/with-wrong-path",
		Channels:      cfg.Exchange.Channels,
		Subscriptions: cfg.Exchange.Subscriptions,
	})

	assert.NotNil(t, err)
}
