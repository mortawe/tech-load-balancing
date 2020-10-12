package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mortawe/chat/internal/chat"
	"github.com/mortawe/chat/internal/errors/apierr"
	"github.com/mortawe/chat/internal/errors/ucerr"
	"github.com/mortawe/chat/internal/models"

	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type ChatHandler struct {
	chatUC chat.UC
}

func NewChatHandler(cUC chat.UC) *ChatHandler {
	return &ChatHandler{chatUC: cUC}
}

func (h *ChatHandler) Register(r *router.Router) {
	r.POST("/chats/add", h.Create)
	r.POST("/chats/get", h.List)
}

func (h *ChatHandler) Create(ctx *fasthttp.RequestCtx) {
	chat := &models.Chat{}
	if err := json.Unmarshal(ctx.PostBody(), &chat); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}
	err := h.chatUC.Create(ctx, chat)
	if errors.Is(err, ucerr.ErrNoUser) {
		ctx.Error(apierr.NoSuchUser, http.StatusBadRequest)
		return
	}
	if errors.Is(err, ucerr.ErrNameAlreadyInUse) {
		ctx.Error(apierr.NameInUse, http.StatusConflict)
		return
	}
	if errors.Is(err, ucerr.ErrUserInChatTwice) {
		ctx.Error(apierr.UserInChatTwice, http.StatusBadRequest)
	}
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		logrus.Error(err)
		return
	}
	ctx.WriteString(fmt.Sprint(chat.ID))
	ctx.SetStatusCode(http.StatusCreated)
}

type ListArgs struct {
	User models.ID `json:"user"`
}

func (h *ChatHandler) List(ctx *fasthttp.RequestCtx) {
	args := &ListArgs{}
	if err := json.Unmarshal(ctx.PostBody(), &args); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}
	chat, err := h.chatUC.GetList(ctx, args.User)
	if errors.Is(err, ucerr.ErrNoUser) {
		ctx.Error(apierr.NoSuchUser, http.StatusBadRequest)
	}
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		logrus.Error(err)
		return
	}
	data, err := json.Marshal(chat)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		logrus.Error(err)
		return
	}
	ctx.Write(data)
	ctx.SetStatusCode(http.StatusOK)
}
