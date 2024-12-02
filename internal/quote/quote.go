package quote

import (
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
)

type Quote struct {
	ID        uuid.UUID
	Message   string
	Person	  string
	CreatedAt time.Time
	IP        net.IP
}

func NewQuote(message string, person string, ip net.IP) (Quote, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return Quote{}, fmt.Errorf("failed to create quote: %w", err)
	}

	return Quote{
		ID:        id,
		Message:   message,
		Person:    person,
		CreatedAt: time.Now(),
		IP:        ip,
	}, nil
}
