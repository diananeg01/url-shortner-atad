package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserId uuid.UUID
	Email  string
}

type Session struct {
	SessionId uuid.UUID
	UserId    uuid.UUID
	ExpiresAt time.Time
}