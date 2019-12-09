package worktime

import (
	"context"

	"github.com/laminh711/reporter/models"
)

// Repository represent the worktime's repository contract
type Repository interface {
	Fetch(ctx context.Context) ([]*models.Worktime, error)
}
