package repository

import (
	"context"

	"github.com/mortawe/chat/internal/errors/dberr"
	"github.com/mortawe/chat/internal/errors/ucerr"
	"github.com/mortawe/chat/internal/models"

	"github.com/jmoiron/sqlx"
)

type MsgRepo struct {
	db *sqlx.DB
}

func NewMsgRepo(db *sqlx.DB) *MsgRepo {
	return &MsgRepo{db: db}
}

const (
	createMsgQ = `INSERT INTO messages (chat_id, author_id, text) 
VALUES (:chat_id, :author_id, :text) RETURNING id`
	getMsgListByChatQ = `SELECT * FROM messages WHERE chat_id=:chat_id ORDER BY created_at`
)

func (r *MsgRepo) Create(ctx context.Context, msg *models.Message) error {
	query, args, err := r.db.BindNamed(createMsgQ, msg)
	if err != nil {
		return err
	}
	err = r.db.GetContext(ctx, &msg.ID, query, args...)
	if dberr.IsForeignKeyViolation(err) {
		return ucerr.ErrUserNotInChat
	}
	return err
}

func (r *MsgRepo) GetList(ctx context.Context, chatID models.ID) ([]models.Message, error) {
	query, args, err := r.db.BindNamed(getMsgListByChatQ, map[string]interface{}{"chat_id": chatID})
	if err != nil {
		return nil, err
	}
	msgList := []models.Message{}
	err = r.db.SelectContext(ctx, &msgList, query, args...)
	return msgList, err
}
