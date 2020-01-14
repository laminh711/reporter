package auth

import (
	"context"

	"github.com/laminh711/reporter/models"
)

type Repository interface {
	Get(context.Context, *models.User) ([]*models.User, error)
}
