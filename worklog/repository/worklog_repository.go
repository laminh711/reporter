package repository

import (
	"context"

	"github.com/laminh711/reporter/models"
	"github.com/laminh711/reporter/worklog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type worklogRepository struct {
	DB         *mongo.Database
	Collection string
}

func NewWorklogRepository(db *mongo.Database, collection string) worklog.Repository {
	return &worklogRepository{db, collection}
}

func (wr *worklogRepository) Fetch(ctx context.Context) ([]*models.Worklog, error) {

	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"finished_at", -1}})
	cursor, err := wr.DB.Collection(wr.Collection).Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	result := []*models.Worklog{}
	worklog := models.Worklog{}
	for cursor.Next(ctx) {
		err := cursor.Decode(&worklog)
		if err != nil {
			return nil, err
		}
		result = append(result, &worklog)
	}

	return result, nil
}

func (wr *worklogRepository) Create(ctx context.Context, worklog models.Worklog) error {
	_, err := wr.DB.Collection(wr.Collection).InsertOne(ctx, worklog)
	// TODO: decide whether if a log is needed here
	if err != nil {
		return err
	}
	return nil
}
