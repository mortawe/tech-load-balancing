package main

import (
	"os"

	. "github.com/mortawe/chat/internal/chat/delivery"
	. "github.com/mortawe/chat/internal/chat/repository"
	. "github.com/mortawe/chat/internal/chat/usecase"
	. "github.com/mortawe/chat/internal/message/delivery"
	. "github.com/mortawe/chat/internal/message/repository"
	. "github.com/mortawe/chat/internal/message/usecase"
	. "github.com/mortawe/chat/internal/user/delivery"
	. "github.com/mortawe/chat/internal/user/repository"
	. "github.com/mortawe/chat/internal/user/usecase"

	"github.com/fasthttp/router"
	"github.com/jmoiron/sqlx"
)

var (
	DBAddr = "localhost:5432"
	DBUser = "user"
	DBPass = "pass"
	DBName = "chat-db"
)

func init() {
	if addr := os.Getenv("DB_ADDR"); addr != "" {
		DBAddr = addr
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		DBName = name
	}
	if user := os.Getenv("DB_USER"); user != "" {
		DBUser = user
	}
	if pass := os.Getenv("DB_PASS"); pass != "" {
		DBPass = pass
	}
}

func register(db *sqlx.DB, r *router.Router) (*UserHandler, *ChatHandler) {
	// user
	uR := NewUserRepo(db)
	uU := NewUserUC(uR)
	uH := NewUserHandler(uU)
	uH.Register(r)
	// chat
	cR := NewChatRepo(db)
	cU := NewChatUC(cR, uR)
	cH := NewChatHandler(cU)
	cH.Register(r)
	// message
	mR := NewMsgRepo(db)
	mU := NewMsgUC(mR, cR)
	mH := NewMsgHandler(mU)
	mH.Register(r)

	return uH, cH
}
