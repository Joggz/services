package usrgrp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Joggz/services/business/data/store/user"
	"github.com/Joggz/services/business/sys/validate"
	"github.com/Joggz/services/business/web/auth"
	"github.com/Joggz/services/foundation/web"
)

type Handlers struct {
	User user.Store
	Auth *auth.Auth
}


func (h Handlers) Create( ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	v, err := web.GetValues(ctx);
	if err != nil {
		return web.NewShutdownError("web value not found in context")
	}	

	var nu user.NewUser
	if err := web.Decode(r, &nu); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	usr, err := h.User.Create(ctx, nu, v.Now)
	if err != nil {
		return fmt.Errorf("user[%+v]: %w", &usr, err)
	}
	return web.Respond(ctx, w, usr, http.StatusCreated)
}

func (h Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page := web.Param(r, "page")
	pageNumber,  err := strconv.Atoi(page)

	if err != nil {
		return validate.NewRequestError(fmt.Errorf("invalid page format [%s]", page), http.StatusBadRequest)
	}

	rows := web.Param(r, "rows")
	rowsPerPage,  err := strconv.Atoi(rows)

	if err != nil {
		return validate.NewRequestError(fmt.Errorf("invalid rows format [%s]", page), http.StatusBadRequest)
	}

	user, err :=	h.User.Query(ctx, pageNumber, rowsPerPage, )

	return web.Respond(ctx, w, user, http.StatusAccepted)
}