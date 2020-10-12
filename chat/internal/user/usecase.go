package user

import (
	"context"

	"github.com/mortawe/chat/internal/models"
)

type UC interface {
	Create(ctx context.Context, user *models.User) error
}
