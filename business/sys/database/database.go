package DB

import (
	"errors"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	q.Set("sslMode", sslMode)
	q.Set("timezone", "UTC")

	// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
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