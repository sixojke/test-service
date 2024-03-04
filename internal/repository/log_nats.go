package repository

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/sixojke/test-service/internal/domain"
)

type LogNats struct {
	nc *nats.Conn
}

func NewLogNats(nc *nats.Conn) *LogNats {
	return &LogNats{nc: nc}
}

func (r *LogNats) Send(msg []*domain.ChangesHistory) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("data serialization: %v", err)
	}

	if err := r.nc.Publish("logs", msgJSON); err != nil {
		return fmt.Errorf("send data: %v", err)
	}

	return nil
}
