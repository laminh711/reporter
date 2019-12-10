package worklog

import (
	"context"

	"github.com/laminh711/reporter/models"
)

type Repository interface {
	Fetch(context.Context) ([]*models.Worklog, error)
	Create(context.Context, models.Worklog) error
}
