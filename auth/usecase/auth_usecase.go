package usecase

import (
	"context"

	"github.com/laminh711/reporter/auth"
	"github.com/laminh711/reporter/models"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	authRepository auth.Repository
}

func NewAuthUsecase(ar auth.Repository) auth.Usecase {
	return &authUsecase{ar}
}

func (au authUsecase) Login(ctx context.Context, user *models.User) error {
	matches, err := au.authRepository.Get(ctx, user)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(matches[0].Password), []byte(user.Password))
	if err != nil {
		return err
	}
	return nil
}
