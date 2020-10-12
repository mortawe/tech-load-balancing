package chat

import (
	"context"

	"github.com/mortawe/chat/internal/models"
)

type Repo interface {
	Create(ctx context.Context, chat *models.Chat) error
	GetList(ctx context.Context, userID models.ID) ([]models.Chat, error)
	Get(ctx context.Context, chatID models.ID) (*models.Chat, error)
}
