package worktime

import (
	"context"

	"github.com/laminh711/reporter/models"
)

// Usecase represent the worktime's use cases
type Usecase interface {
	Fetch(ctx context.Context) ([]*models.Worktime, error)
}
