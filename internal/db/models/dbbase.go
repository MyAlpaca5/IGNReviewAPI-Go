package models

import "time"

// dbBase struct includes data that is unique for database but may be not valuable for external usage, for example record created time
type dbBase struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
}
