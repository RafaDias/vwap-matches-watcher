package coinbase

import (
	"encoding/json"
	"fmt"
	"github.com/rafadias/crypto-watcher/internal/application/providers/exchange"
	"github.com/rafadias/crypto-watcher/internal/domain"
	"golang.org/x/net/websocket"
	"log"
	"strconv"
)

const (
	host  = "http://localhost"
)

type service struct {
	client *websocket.Conn
	log *log.Logger
	subscriptions []string
	channels []string
}

func (s *service) GetSubscriptions() []string {
	return s.subscriptions
}

func (s *service) ListenTransactions(c chan domain.Transaction) error {
	if err := s.subscribe(); err != nil {
		return err
	}

	for {
		var msg = make([]byte, 512) // 276B
		bodyLength, err := s.client.Read(msg)
		if err != nil {
			s.log.Println(fmt.Errorf("fail to read incoming message: %v", err))
			return err
		}

		var response Response
		if err = json.Unmarshal(msg[:bodyLength], &response); err != nil {
			log.Fatal("error trying to unmarshal response")
		}
		if response.Type != Subscriptions {
			trx := translateResponseToDomain(response)
			c <- trx
		}
	}
}

func New(config exchange.Config) (exchange.Service, error) {
	conn, err := websocket.Dial(config.BaseUrl, "", host)
	if err != nil {
		log.Println("error trying to connect to coinbase")
		return nil, err
	}
	log.Println("init exchange with config: ", config)

	return &service{
		client: conn,
		channels: config.Channels,
		subscriptions: config.Subscriptions,
	}, nil
}

func translateResponseToDomain(response Response) domain.Transaction {
	price, err := strconv.ParseFloat(response.Price, 64)
	if err != nil {
		log.Fatal("err trying to convert price")
	}
	return domain.Transaction{
		ProductId: response.ProductID,
		Price:     price,
		Size: response.Size,
		Time: response.Time,
	}
}

func (s *service) subscribe() error {
	subscriptionMessage := &Message{
		Type:       Subscribe,
		ProductIDs: s.subscriptions,
		Channels:   s.channels,
	}
	msg, err := json.Marshal(subscriptionMessage)
	if err != nil {
		log.Fatal()
	}
	if _, err = s.client.Write(msg); err != nil {
		return err
	}
	return nil
}
