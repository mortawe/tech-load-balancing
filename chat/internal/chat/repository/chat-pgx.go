package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mortawe/chat/internal/errors/dberr"
	"github.com/mortawe/chat/internal/errors/ucerr"
	"github.com/mortawe/chat/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ChatRepo struct {
	db *sqlx.DB
}

func NewChatRepo(db *sqlx.DB) *ChatRepo {
	return &ChatRepo{db: db}
}

const (
	createChatQ      = "INSERT INTO chats (name) VALUES (:name) RETURNING id"
	insertChatUsersQ = "INSERT INTO chat_users (chat_id, user_id) VALUES "
	getChatsByUserQ  = `SELECT chats.id, chats.name, chats.created_at, array_users.users AS users
FROM chat_users
         JOIN chats ON chat_users.chat_id = chats.id AND chat_users.user_id = $1
         LEFT JOIN (SELECT max(created_at) AS created_at, chat_id
               FROM messages
               GROUP BY chat_id
            )  AS last_modify ON chat_users.chat_id = last_modify.chat_id
         JOIN (SELECT chat_id, array_agg(user_id) AS users
               FROM chat_users
               GROUP BY chat_id) AS array_users ON array_users.chat_id = chat_users.chat_id
ORDER BY last_modify.created_at DESC NULLS LAST`
	getChatQ = `SELECT * FROM chats WHERE id = :chat_id`
)

func buildQueryInsertValues(chatID models.ID, ids []models.ID) string {
	res := ""
	for _, id := range ids {
		res = fmt.Sprint(res, "(", chatID, ",", id, "),")
	}
	if len(res) > 0 {
		return res[:len(res)-1]
	}
	return ""
}

func (r *ChatRepo) Create(ctx context.Context, chat *models.Chat) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query, args, err := tx.BindNamed(createChatQ, &chat)
	if err != nil {
		return err
	}

	if err := tx.GetContext(ctx, &chat.ID, query, args...); err != nil {
		if dberr.IsUniqueViolationErr(err) {
			return ucerr.ErrNameAlreadyInUse
		}
		return err
	}
	if len(chat.Users) == 0 {
		return nil
	}
	values := buildQueryInsertValues(chat.ID, chat.Users)
	if _, err := tx.ExecContext(ctx, insertChatUsersQ+values); err != nil {
		if dberr.IsUniqueViolationErr(err) {
			return ucerr.ErrUserInChatTwice
		}
		if dberr.IsForeignKeyViolation(err) {
			return ucerr.ErrNoUser
		}
		return err
	}

	return tx.Commit()
}

func (r *ChatRepo) GetList(ctx context.Context, userID models.ID) ([]models.Chat, error) {
	chats := []models.Chat{}
	rows, err := r.db.QueryContext(ctx, getChatsByUserQ, &userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		chat := models.Chat{}
		users := []int64{}
		if err := rows.Scan(&chat.ID, &chat.Name, &chat.CreatedAt, pq.Array(&users)); err != nil {
			return nil, err
		}
		chat.Users = models.CastInt64ArrToIdArr(users)
		chats = append(chats, chat)
	}
	return chats, err
}

func (r *ChatRepo) Get(ctx context.Context, chatID models.ID) (*models.Chat, error) {
	query, args, err := r.db.BindNamed(getChatQ, map[string]interface{}{"chat_id": chatID})
	if err != nil {
		return nil, err
	}
	chat := models.Chat{}
	err = r.db.GetContext(ctx, &chat, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ucerr.ErrNoChat
	}
	return &chat, err
}
