package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Log Log
type Log struct {
	ID         *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name       string              `json:"name"`
	Productive bool                `json:"productive"`
	StartAt    time.Time           `json:"start_at"`
}
