package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mortawe/chat/internal/errors/apierr"
	"github.com/mortawe/chat/internal/errors/ucerr"
	"github.com/mortawe/chat/internal/models"
	"github.com/mortawe/chat/internal/user"

	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type UserHandler struct {
	userUC user.UC
}

func NewUserHandler(uUC user.UC) *UserHandler {
	return &UserHandler{userUC: uUC}
}

func (h *UserHandler) Register(r *router.Router) {
	r.POST("/users/add", h.Create)
}

func (h *UserHandler) Create(ctx *fasthttp.RequestCtx) {
	user := &models.User{}
	if err := json.Unmarshal(ctx.PostBody(), &user); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}
	err := h.userUC.Create(ctx, user)
	if errors.Is(err, ucerr.ErrNameAlreadyInUse) {
		ctx.Error(apierr.NameInUse, http.StatusConflict)
		return
	}
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		logrus.Error(err)
		return
	}
	ctx.WriteString(fmt.Sprint(user.ID))
	ctx.SetStatusCode(http.StatusCreated)
}
