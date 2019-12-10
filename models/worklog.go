package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Worklog Worklog
type Worklog struct {
	ID         *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name       string              `json:"name" bson:"name"`
	Productive bool                `json:"productive" bson:"productive"`
	FinishedAt time.Time           `json:"finished_at" bson:"finished_at"`
}
