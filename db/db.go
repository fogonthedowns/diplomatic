package db

import (
	"context"

	"diplomacy/model"
)

type Crud interface {
	// Fetch(ctx context.Context, num int64) ([]*models.Post, error)
	// GetByID(ctx context.Context, id int64) (*models.Post, error)
	Create(ctx context.Context, p *model.Game) (int64, error)
	// Update(ctx context.Context, p *models.Post) (*models.Post, error)
	// Delete(ctx context.Context, id int64) (bool, error)
}