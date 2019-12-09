package repository

import (
	"context"

	"github.com/laminh711/reporter/models"
	"github.com/laminh711/reporter/worktime"

	"github.com/k0kubun/pp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mongoWorktimeRepository struct {
	DB         *mongo.Database
	Collection string
}

// NewMongoWorktimeRepository create new mongodb connection for worktime
func NewMongoWorktimeRepository(Database *mongo.Database, Collection string) worktime.Repository {
	return &mongoWorktimeRepository{Database, Collection}
}

func (m *mongoWorktimeRepository) Fetch(ctx context.Context) ([]*models.Worktime, error) {
	// collection := m.DB.Collection(m.Collection)
	filter := bson.M{}
	// findOptions := options.Find()
	// findOptions.SetSort(bson.D{{"left_end", 1}})
	cursor, err := m.DB.Collection(m.Collection).Find(ctx, filter)
	if err != nil {
		pp.Println("34 mongo worktime")
		return nil, err
	}

	// result := []*models.Worktime{}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		worktimeRecord := models.Worktime{}
		err := cursor.Decode(&worktimeRecord)
		if err != nil {
			return nil, err
		}
		// result = append(result, &worktimeRecord)
	}
	return nil, nil
	// return result, nil
}
