package message

import (
	"context"

	"github.com/mortawe/chat/internal/models"
)

type UC interface {
	Create(ctx context.Context, chat *models.Message) error
	GetList(ctx context.Context, chatID models.ID) ([]models.Message, error)
}
