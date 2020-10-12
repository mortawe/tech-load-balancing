package models

import (
	"time"
)

type ID int64 // int64 or string

type User struct {
	ID        ID        `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Chat struct {
	ID        ID        `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Users     []ID      `json:"users" db:"users"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Message struct {
	ID        ID        `json:"id" db:"id"`
	ChatID    ID        `json:"chat" db:"chat_id"`
	AuthorID  ID        `json:"author" db:"author_id"`
	Text      string    `json:"text" db:"text"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func CastInt64ArrToIdArr(arr []int64) []ID {
	result := make([]ID, len(arr))
	for i, e := range arr {
		temp := ID(e)
		result[i] = temp
	}
	return result
}
