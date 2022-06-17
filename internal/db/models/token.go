package models

import (
	"time"
)

type Token struct {
	TokenHash []byte
	UserID    int
	Expiry    time.Time
	Role      int
}
