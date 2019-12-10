package usecase

import (
	"context"

	"github.com/laminh711/reporter/models"
	"github.com/laminh711/reporter/worklog"
)

type worklogUsecase struct {
	worklogRepository worklog.Repository
}

func NewWorklogUsecase(wr worklog.Repository) worklog.Usecase {
	return &worklogUsecase{wr}
}

func (wu *worklogUsecase) Fetch(ctx context.Context) ([]*models.Worklog, error) {
	result, err := wu.worklogRepository.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (wu *worklogUsecase) Create(ctx context.Context, worklog models.Worklog) error {
	err := wu.worklogRepository.Create(ctx, worklog)
	if err != nil {
		return err
	}
	return nil
}
