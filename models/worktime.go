package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Worktime model
type Worktime struct {
	ID            *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	LeftEnd       time.Time           `json:"left_end" bson:"left_end"`
	RightEnd      time.Time           `json:"right_end" bson:"right_end"`
	Interval      time.Time           `json:"interval" bson:"interval"`
	BreakDuration time.Time           `json:"break_duration" bson:"break_duration"`
}
