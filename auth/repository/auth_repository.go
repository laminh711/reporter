package repository

import (
	"context"

	"github.com/laminh711/reporter/auth"
	"github.com/laminh711/reporter/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository struct {
	DB         *mongo.Database
	Collection string
}

func NewAuthRepository(db *mongo.Database, collection string) auth.Repository {
	return AuthRepository{
		db,
		collection,
	}
}

func (ar AuthRepository) Get(ctx context.Context, user *models.User) ([]*models.User, error) {
	filter := bson.M{
		"username": user.Username,
	}

	cursor, err := ar.DB.Collection(ar.Collection).Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	result := []*models.User{}
	bindUser := models.User{}
	for cursor.Next(ctx) {
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		result = append(result, &bindUser)
	}

	return result, nil
}
