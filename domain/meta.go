package domain

import "time"

// Meta represents standard postgres fields to be embedded in database models.
type Meta struct {
	ID        uint64     `json:"id,omitempty" db:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}
