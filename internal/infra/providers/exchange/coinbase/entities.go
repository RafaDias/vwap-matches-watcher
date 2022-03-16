package coinbase

import "time"

const (
	Subscribe     string = "subscribe"
	Subscriptions string = "subscriptions"
	Match         string = "match"
)

type Message struct {
	Type       string   `json:"type"`
	ProductIDs []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

type Response struct {
	Type         string    `json:"type"`
	TradeID      int       `json:"trade_id"`
	MakerOrderID string    `json:"maker_order_id"`
	TakerOrderID string    `json:"taker_order_id"`
	Side         string    `json:"side"`
	Size         string    `json:"size"`
	Price        string    `json:"price"`
	ProductID    string    `json:"product_id"`
	Sequence     int64     `json:"sequence"`
	Time         time.Time `json:"time"`
}
