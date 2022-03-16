package coinbase

import (
	"github.com/joho/godotenv"
	"github.com/rafadias/crypto-watcher/internal/application/providers/exchange"
	"github.com/rafadias/crypto-watcher/internal/config"
	"github.com/rafadias/crypto-watcher/internal/domain"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

const filePath = "../../../../../.env"

func TestService_Connect_With_Correct_Values(t *testing.T) {
	if err := godotenv.Load(filePath); err != nil {
		log.Print("No .env file found")
	}

	cfg := config.New()
	log := log.New(os.Stdout, "testing", log.Lshortfile)

	svc, err := New(exchange.Config{
		BaseURL:       cfg.Exchange.BaseURL,
		Channels:      cfg.Exchange.Channels,
		Subscriptions: cfg.Exchange.Subscriptions,
	}, log, true)

	assert.Nil(t, err)
	assert.Equal(t, cfg.Exchange.Subscriptions, svc.GetSubscriptions())
}

func TestService_WrongConnectionMustReturnAnError(t *testing.T) {
	if err := godotenv.Load(filePath); err != nil {
		log.Print("No .env file found")
	}

	cfg := config.New()
	log := log.New(os.Stdout, "testing", log.Lshortfile)

	_, err := New(exchange.Config{
		BaseURL:       "wss:/wrong-host/with-wrong-path",
		Channels:      cfg.Exchange.Channels,
		Subscriptions: cfg.Exchange.Subscriptions,
	}, log, true)

	assert.NotNil(t, err)
}

func TestService_GetChannels(t *testing.T) {
	channels := []string{"first-channel"}
	log := log.New(os.Stdout, "testing", log.Lshortfile)

	svc, _ := New(exchange.Config{
		Channels: channels,
	}, log, false)
	assert.Equal(t, channels, svc.GetChannels())
}

func TestService_GetSubscriptions(t *testing.T) {
	subscriptions := []string{"subs1", "subs2"}
	log := log.New(os.Stdout, "testing", log.Lshortfile)

	svc, _ := New(exchange.Config{
		Subscriptions: subscriptions,
	}, log, false)
	assert.Equal(t, subscriptions, svc.GetSubscriptions())
}

func TestTranslateResponseToDomain(t *testing.T) {
	resp := Response{
		ProductID:    "ETH-USD",
		Price:        "2620.27",
		MakerOrderID: "958e3d82-1845-4e12-aff0-e2284e77b8d2",
		TakerOrderID: "cb209aeb-ea08-410b-9d9d-f87296ed82b5",
		Size:         "0.00342402",
		Side:         "sell",
		Time:         time.Now(),
	}
	txn, _ := translateResponseToDomain(resp)

	price, err := strconv.ParseFloat(resp.Price, 64)
	if err != nil {
		t.Fatal("err trying to convert price")
	}
	assert.Equal(t, price, txn.Price)

	size, err := strconv.ParseFloat(resp.Size, 64)
	if err != nil {
		t.Fatal("err trying to convert price")
	}
	assert.Equal(t, size, txn.Size)

	assert.Equal(t, txn.ProductID, resp.ProductID)
	assert.Equal(t, txn.Time, resp.Time)

}

func TestWrongPriceShouldReturnAnError(t *testing.T) {
	resp := Response{Price: "bla bla"}
	_, err := translateResponseToDomain(resp)
	assert.ErrorIs(t, err, domain.ErrCannotConvertPrice)
}

func TestWrongSizeShouldReturnAnError(t *testing.T) {
	resp := Response{Size: "bla bla", Price: "1.0"}
	_, err := translateResponseToDomain(resp)
	assert.ErrorIs(t, err, domain.ErrCannotConvertSize)
}

func TestSubscribe(t *testing.T) {
	subscriptions := []string{"subs1", "subs2"}
	channels := []string{"first-channel"}
	svc := &service{
		subscriptions: subscriptions,
		channels:      channels,
		dial:          false,
	}
	msg, err := svc.subscribe()
	assert.Nil(t, err)
	assert.Equal(t, Subscribe, msg.Type)
	assert.Equal(t, subscriptions, msg.ProductIDs)
	assert.Equal(t, channels, msg.Channels)
}
