// Package user provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Joggz/services/business/data/store/user"
	"github.com/Joggz/services/business/web/auth"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Core manages the set of APIs for user access.
type Core struct {
	log *zap.SugaredLogger
	user user.Store
}

// NewCore constructs a core for user api access.
func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		log: log,
		user: user.NewStore(log, sqlxDB),
	}
}


func (c Core) Create(ctx context.Context, nu user.NewUser, now time.Time ) (user.User, error) {

		// PERFORM PRE BUSINESS OPERATION
	usr, err :=	c.user.Create(ctx, nu, now )
	if err != nil {
		return user.User{}, fmt.Errorf("create :%w",  err)
	}
		// PERFORM POST BUSINESS OPERATION
	return usr, nil
}

func (c Core) Update(ctx context.Context, claim auth.Claims,  userID string, uu user.UpdateUser, now time.Time) error {

	// PERFORM PRE BUSINESS OPERATION
	 err :=	c.user.Update(ctx, claim, userID, uu, now )
		if err != nil {
			return fmt.Errorf("update :%w",  err)
		}
	// PERFORM POST BUSINESS OPERATION
		return  nil
}


func (c Core) Delete(ctx context.Context, claim auth.Claims, userID string ) (error) {

	// PERFORM PRE BUSINESS OPERATION
		 err :=	c.user.Delete(ctx, claim, userID )
		if err != nil {
			return fmt.Errorf("delete :%w",  err)
		}
	// PERFORM POST BUSINESS OPERATION
		return  nil
}


func (c Core) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]user.User, error) {

	// PERFORM PRE BUSINESS OPERATION
		usr, err :=	c.user.Query(ctx, pageNumber, rowsPerPage )
		if err != nil {
			return []user.User{}, fmt.Errorf("query :%w",  err)
		}
	// PERFORM POST BUSINESS OPERATION
		return usr, nil
}


func (c Core) QueryByID(ctx context.Context, claim auth.Claims, userID string ) (user.User, error) {

	// PERFORM PRE BUSINESS OPERATION
		usr, err :=	c.user.QueryByID(ctx, claim, userID )
		if err != nil {
			return user.User{}, fmt.Errorf("create :%w",  err)
		}
	// PERFORM POST BUSINESS OPERATION
		return usr, nil
}


func (c Core) QueryByEmail(ctx context.Context, claim auth.Claims, email string) (user.User, error) {
		// PERFORM PRE BUSINESS OPERATION
		usr, err :=	c.user.QueryByEmail(ctx, claim, email )
		if err != nil {
			return user.User{}, fmt.Errorf("create :%w",  err)
		}
			// PERFORM POST BUSINESS OPERATION
		return usr, nil
}


func (c Core) Authenticate(ctx context.Context, now time.Time, email string, password string) (auth.Claims, error){
	// PERFORM PRE BUSINESS OPERATION
	claim, err :=	c.user.Authenticate(ctx, now, email, password )
	if err != nil {
		return auth.Claims{}, fmt.Errorf("create :%w",  err)
	}
		// PERFORM POST BUSINESS OPERATION
	return claim, nil
}