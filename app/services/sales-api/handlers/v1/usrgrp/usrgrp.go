package usrgrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Joggz/services/business/data/store/user"
	"github.com/Joggz/services/business/sys/database"
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

func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claim, err := auth.GetClaims(ctx)
	if err != nil {
		return errors.New("cant find clains in context")
	}

	id := web.Param(r, "id")
	user, err := h.User.QueryByID(ctx, claim, id)
	if err != nil {
		switch{
		case errors.Is(err, database.ErrInvalidID):
			return validate.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, database.ErrNotFound):
			return validate.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, w, user, http.StatusAccepted)
}

func (h Handlers) Update( ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value not found in context")
	}	

	claim, err := auth.GetClaims(ctx)	
	if err != nil {
		return errors.New("cant find clains in context")
	}

	id := web.Param(r, "id")

	var upd user.UpdateUser
	if err := web.Decode(r, &upd); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	if err := h.User.Update(ctx, claim, id, upd, v.Now); err != nil {
		switch{
		case errors.Is(err, database.ErrInvalidID):
			return validate.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, database.ErrNotFound):
			return validate.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", id, err)
		}
	}
	return web.Respond(ctx, w, upd, http.StatusNoContent)
}

func (h Handlers) QueryByEmail(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claim, err := auth.GetClaims(ctx)
	if err != nil {
		return errors.New("cant find clains in context")
	}
	email := web.Param(r, "email")
	user, err := h.User.QueryByEmail(ctx, claim, email)
	if err != nil {
		switch  {
		case errors.Is(err, database.ErrInvalidID):
			return validate.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, database.ErrNotFound):
			return validate.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", email, err)
		}
			
		}
		return web.Respond(ctx, w, user, http.StatusNoContent)
}

func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claim, err := auth.GetClaims(ctx)
	if err != nil {
		return errors.New("cant find clains in context")
	}

	id := web.Param(r, "id")
	derr := h.User.Delete(ctx, claim, id);
	if derr != nil {
		switch {
		case errors.Is(err, user.ErrInvalidID):
			return validate.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("ID[%s]: %w", id, err)
		}
	}
		return web.Respond(ctx, w, user.User{}, http.StatusAccepted)
}


// Token provides an API token for the authenticated user.
func (h Handlers) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	email, pass, ok := r.BasicAuth()
	fmt.Println("email && pass in usrgrp file", email, pass)
	if !ok {
		err := errors.New("must provide email and password in Basic auth")
		return validate.NewRequestError(err, http.StatusUnauthorized)
	}

	claims, err := h.User.Authenticate(ctx, v.Now, email, pass)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return validate.NewRequestError(err, http.StatusNotFound)
		case errors.Is(err, user.ErrAuthenticationFailure):
			return validate.NewRequestError(err, http.StatusUnauthorized)
		default:
			return fmt.Errorf("authenticating: %w", err)
		}
	}

	var tkn struct {
		Token string `json:"token"`
	}
	tkn.Token, err = h.Auth.GenerateToken(claims)
	if err != nil {
		return fmt.Errorf("generating token: %w", err)
	}

	return web.Respond(ctx, w, tkn, http.StatusOK)
}
