package auth

import (
	"context"

	"github.com/laminh711/reporter/models"
)

type Usecase interface {
	Login(context.Context, *models.User) error
}
