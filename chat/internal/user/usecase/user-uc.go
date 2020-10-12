package usecase

import (
	"context"

	"github.com/mortawe/chat/internal/models"
	"github.com/mortawe/chat/internal/user"
)

type UserUC struct {
	userRepo user.Repo
}

func NewUserUC(uR user.Repo) *UserUC {
	return &UserUC{userRepo: uR}
}

func (u *UserUC) Create(ctx context.Context, user *models.User) error {
	return u.userRepo.Create(ctx, user)
}
