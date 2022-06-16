package database

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)


var (
	ErrNotFound  = errors.New("not found")
	ErrInvalidID  = errors.New("ID is not in its proper form")
	ErrAuthenticationFailure = errors.New("authentication failed")
	ErrForbidden = errors.New("attempted action is not allowed")
)

// Config is the required properties to use the database.
type Config struct {
	User         string
	Password     string 
	Host         string
	Name         string
	MaxIdleConns int
	MaxOpenConns int
	DisableTLS   bool
}


// Open knows how to open a database connection based on the configuration.
 func Open (cfg  Config) (*sqlx.DB, error){
	
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	
	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "UTC")

	// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
	fmt.Println("SSLMODE %w", sslMode, q.Encode())
	u := url.URL{
		Scheme:      "postgres",
		Opaque:      "",
		User:        url.UserPassword(cfg.User, cfg.Password),
		Host:        cfg.Host,
		Path:        cfg.Name,
		RawQuery:    q.Encode(),	
	};

	// u =  "postgres://usernamepassword@host/........."

	 db, err := sqlx.Open("postgres", u.String())
	 if err != nil {
		 return nil, err
	 }


	 db.SetMaxIdleConns(cfg.MaxIdleConns)
	 db.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
 }

 // StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sqlx.DB, log *zap.SugaredLogger) error  {
	

	var pingError error

	for attempts := 1; ; attempts++ {
		pingError =  db.Ping()
		if pingError == nil {
			break;	
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		// Make sure we didn't timeout or be cancelled.
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	

	// Make sure we didn't timeout or be cancelled.
	if ctx.Err() != nil {
		return ctx.Err()
	}

		// Run a simple query to determine connectivity. Running this query forces a
	// round trip through the database.
	const q = `SELECT true`
	// var tmp bool;
	
	// return db.QueryRowContext(ctx, q).Scan(&tmp)
	return nil
}

