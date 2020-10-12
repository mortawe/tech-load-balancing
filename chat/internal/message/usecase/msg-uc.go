package usecase

import (
	"context"

	"github.com/mortawe/chat/internal/chat"
	"github.com/mortawe/chat/internal/message"
	"github.com/mortawe/chat/internal/models"
)

type MsgUC struct {
	msgRepo  message.Repo
	chatRepo chat.Repo
}

func NewMsgUC(mR message.Repo, cR chat.Repo) *MsgUC {
	return &MsgUC{msgRepo: mR, chatRepo: cR}
}

func (u *MsgUC) Create(ctx context.Context, msg *models.Message) error {
	return u.msgRepo.Create(ctx, msg)
}

func (u *MsgUC) GetList(ctx context.Context, chatID models.ID) ([]models.Message, error) {
	if _, err := u.chatRepo.Get(ctx, chatID); err != nil {
		return nil, err
	}
	return u.msgRepo.GetList(ctx, chatID)
}
