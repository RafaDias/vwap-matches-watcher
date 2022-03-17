// Package websocketserver implements a server for integration testing purposes
package websocketserver

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
)

var upgrader = websocket.Upgrader{}

type hub struct {
	connections []*websocket.Conn
}

type server struct {
	svr *httptest.Server
	hub *hub
}

func (s *server) URL() string {
	return s.svr.URL
}

func New() *server {
	h := new(hub)
	return &server{
		hub: h,
		svr: httptest.NewServer(http.HandlerFunc(h.echo)),
	}
}

func (h *hub) echo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entrou")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	h.connections = append(h.connections, c)

	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		for _, conn := range h.connections {
			err = conn.WriteMessage(mt, message)
			if err != nil {
				break
			}
		}

	}
}
