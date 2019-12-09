package usecase

import (
	"context"
	"time"

	"github.com/laminh711/reporter/models"
	"github.com/laminh711/reporter/worktime"
)

type worktimeUsecase struct {
	worktimeRepo   worktime.Repository
	contextTimeout time.Duration
}

// NewWorktimeUsecase create Worktime entity use cases
func NewWorktimeUsecase(wr worktime.Repository, timeout time.Duration) worktime.Usecase {
	return &worktimeUsecase{
		worktimeRepo:   wr,
		contextTimeout: timeout,
	}
}

func (wu *worktimeUsecase) Fetch(ctx context.Context) ([]*models.Worktime, error) {
	result, err := wu.worktimeRepo.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}
