package coinbase

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rafadias/crypto-watcher/internal/application/providers/exchange"
	"github.com/rafadias/crypto-watcher/internal/config"
	"github.com/rafadias/crypto-watcher/internal/domain"
	"github.com/rafadias/crypto-watcher/internal/infra/websocketserver"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
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

func TestService_ListenTransactions(t *testing.T) {
	svr := websocketserver.New()
	log := log.New(os.Stdout, "testing", log.Lshortfile)

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(svr.URL(), "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	coinbase, err := New(exchange.Config{
		BaseURL:       u,
		Channels:      []string{"channels"},
		Subscriptions: []string{"subs"},
	}, log, true)
	if err != nil {
		return
	}

	transaction := make(chan domain.Transaction)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := coinbase.ListenTransactions(transaction)
		if err != nil {
			log.Print("Error trying to listen")
		}
	}()

	now := time.Now()
	response := Response{Size: "1.0", Price: "3.0", Type: Match, Time: now}
	expectedTxn := domain.Transaction{Size: 1.0, Price: 3.0, Time: now}
	responseJson, err := json.Marshal(response)
	err = ws.WriteMessage(websocket.TextMessage, responseJson)
	if err != nil {
		t.Fatalf("%v", err)
	}

	txn, ok := <-transaction
	if !ok {
		t.Fatalf("msg was not received")
	}
	assert.Equal(t, expectedTxn.Price, txn.Price)
	assert.Equal(t, expectedTxn.Size, txn.Size)
	assert.True(t, expectedTxn.Time.Equal(txn.Time))
}
