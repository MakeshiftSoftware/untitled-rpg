package domain

import "time"

// Meta represents standard database fields. This struct can be
// embedded within another struct to automatically populate
// these fields when scanning a query result.
type Meta struct {
	ID        uint64     `json:"id,omitempty" db:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}
