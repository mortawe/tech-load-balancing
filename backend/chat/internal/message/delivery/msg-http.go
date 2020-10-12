package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mortawe/chat/internal/errors/apierr"
	"github.com/mortawe/chat/internal/errors/ucerr"
	"github.com/mortawe/chat/internal/message"
	"github.com/mortawe/chat/internal/models"

	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type MsgHandler struct {
	msgUC message.UC
}

func NewMsgHandler(mU message.UC) *MsgHandler {
	return &MsgHandler{msgUC: mU}
}

func (h *MsgHandler) Register(r *router.Router) {
	r.POST("/messages/add", h.Create)
	r.POST("/messages/get", h.List)
}

func (h *MsgHandler) Create(ctx *fasthttp.RequestCtx) {
	msg := &models.Message{}
	if err := json.Unmarshal(ctx.PostBody(), &msg); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}
	err := h.msgUC.Create(ctx, msg)
	if errors.Is(err, ucerr.ErrUserNotInChat) {
		ctx.Error(apierr.UserNotInChat, http.StatusForbidden)
		return
	}
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		logrus.Error(err)
		return
	}
	ctx.WriteString(fmt.Sprint(msg.ID))
	ctx.SetStatusCode(http.StatusCreated)
}

type ListMsgArgs struct {
	ChatID models.ID `json:"chat"`
}

func (h *MsgHandler) List(ctx *fasthttp.RequestCtx) {
	args := &ListMsgArgs{}
	if err := json.Unmarshal(ctx.PostBody(), args); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}
	msgList, err := h.msgUC.GetList(ctx, args.ChatID)
	if errors.Is(err, ucerr.ErrNoChat) {
		ctx.Error(apierr.NoSuchChat, http.StatusBadRequest)
		return
	}
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		logrus.Error(err)
		return
	}
	data, err := json.Marshal(msgList)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		logrus.Error(err)
		return
	}
	ctx.Write(data)
	ctx.SetStatusCode(http.StatusOK)
}
