// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type Guest struct {
	ID        uuid.UUID
	Message   string
	Ip        net.IP
	CreatedAt time.Time
	UpdatedAt time.Time
}