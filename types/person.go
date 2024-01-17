package types

import (
	"time"
)

type Person struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string    `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string    `json:"lastname,omitempty" bson:"lastname,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
