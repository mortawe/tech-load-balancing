package usecase

import (
	"context"

	"github.com/mortawe/chat/internal/chat"
	"github.com/mortawe/chat/internal/models"
	"github.com/mortawe/chat/internal/user"
)

type ChatUC struct {
	chatRepo chat.Repo
	userRepo user.Repo
}

func NewChatUC(cR chat.Repo, uR user.Repo) *ChatUC {
	return &ChatUC{chatRepo: cR, userRepo: uR}
}

func (u *ChatUC) Create(ctx context.Context, chat *models.Chat) error {
	return u.chatRepo.Create(ctx, chat)
}

func (u *ChatUC) GetList(ctx context.Context, userID models.ID) ([]models.Chat, error) {
	if _, err := u.userRepo.Get(ctx, userID); err != nil {
		return nil, err
	}
	return u.chatRepo.GetList(ctx, userID)
}
