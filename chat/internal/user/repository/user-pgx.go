package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mortawe/chat/internal/errors/dberr"
	"github.com/mortawe/chat/internal/errors/ucerr"
	"github.com/mortawe/chat/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

const (
	createUserQ = "INSERT INTO users (username) VALUES (:username) RETURNING id"
	getUserQ    = `SELECT * FROM users WHERE id = :id`
)

func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	query, args, err := r.db.BindNamed(createUserQ, &user)
	if err != nil {
		return err
	}
	err = r.db.GetContext(ctx, &user.ID, query, args...)
	if dberr.IsUniqueViolationErr(err) {
		return ucerr.ErrNameAlreadyInUse
	}
	return err
}

func (r *UserRepo) Get(ctx context.Context, userID models.ID) (*models.User, error) {
	user := &models.User{}
	query, args, err := r.db.BindNamed(getUserQ, map[string]interface{}{"id": userID})
	if err != nil {
		return nil, err
	}
	err = r.db.GetContext(ctx, user, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ucerr.ErrNoUser
	}
	return user, err
}
