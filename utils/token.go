package utils

import (
	"paysyncevets/models"
	"time"
)

type Payload struct {
	User models.User `json:"user"`
}

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}