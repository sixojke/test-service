package repository

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

type Message struct {
	Name string `json:"name"`
}

func NewNatsClient() (*nats.Conn, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, fmt.Errorf("connetction: %v", err)
	}

	return nc, nil
}
