// Package coinbase should implement exchange APi and allow access to coinbase platform
package coinbase

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/rafadias/crypto-watcher/internal/application/providers/exchange"
	"github.com/rafadias/crypto-watcher/internal/domain"
	"golang.org/x/net/websocket"
)

const (
	host = "http://localhost"
)

type service struct {
	client        *websocket.Conn
	log           *log.Logger
	subscriptions []string
	channels      []string
	dial          bool
}

func (s *service) GetChannels() []string {
	return s.channels
}

func (s *service) Close() error {
	return s.client.Close()
}

func (s *service) GetSubscriptions() []string {
	return s.subscriptions
}

func (s *service) ListenTransactions(c chan domain.Transaction) error {
	defer close(c)
	defer s.client.Close()

	if _, err := s.subscribe(); err != nil {
		return err
	}

	for {
		_, msg, err := s.client.ReadMessage()
		if err != nil {
			s.log.Println(fmt.Errorf("fail to read incoming message: %v", err))
			return err
		}

		var response Response
		if err = json.Unmarshal(msg, &response); err != nil {
			s.log.Println("error trying to unmarshal response:  ", string(msg))
		}
		if response.Type != Subscriptions {
			trx, err := translateResponseToDomain(response)
			if err != nil {
				s.log.Println("err trying to translate response", err)
			}
			c <- trx
		}
	}
}

func New(config exchange.Config, dial bool) (exchange.Service, error) {
	svc := &service{
		channels:      config.Channels,
		subscriptions: config.Subscriptions,
		dial:          dial,
	}
	if dial {
		conn, _, err := websocket.DefaultDialer.Dial(config.BaseURL, nil)
		if err != nil {
			log.Println("error trying to connect to coinbase")
			return nil, err
		}
		svc.client = conn
	}

	log.Println("init exchange with config: ", config)

	return svc, nil
}

func translateResponseToDomain(response Response) (domain.Transaction, error) {
	price, err := strconv.ParseFloat(response.Price, 64)
	if err != nil {
		return domain.Transaction{}, domain.ErrCannotConvertPrice
	}

	size, err := strconv.ParseFloat(response.Size, 64)
	if err != nil {
		return domain.Transaction{}, domain.ErrCannotConvertSize
	}

	return domain.Transaction{
		ProductID: response.ProductID,
		Price:     price,
		Size:      size,
		Time:      response.Time,
	}, nil
}

func (s *service) subscribe() (*Message, error) {
	subscriptionMessage := &Message{
		Type:       Subscribe,
		ProductIDs: s.subscriptions,
		Channels:   s.channels,
	}
	if !s.dial {
		return subscriptionMessage, nil
	}

	msg, err := json.Marshal(subscriptionMessage)
	if err != nil {
		s.log.Fatal("err trying to write subscribe message")
	}
	if err = s.client.WriteMessage(websocket.TextMessage, msg); err != nil {
		return nil, err
	}

	return nil, nil
}
