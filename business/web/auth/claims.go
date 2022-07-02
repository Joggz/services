package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

const (
	RoleAdmin = "ADMIN"
	RoleUser = "USER"
)

type Claims struct {
	jwt.RegisteredClaims;
	Role []string `json:"roles"`
}

func(c Claims) Authorized(role ...string) bool {
	for _, has := range c.Role {
		for _, want := range role {
			if has == want {
				return true
			}
		}
	}
	return false
}


// ctxKey represents the type of value for the context key.
type ctxkey int

// key is used to store/retrieve a Claims value from a context.Context.
const key ctxkey = 1


// SetClaims stores the claims in the context.
func SetClaims(ctx context.Context, claim Claims) context.Context {
	return context.WithValue(ctx, key, claim)
}

// GetClaims returns the claims from the context.
func GetClaims(ctx context.Context) (Claims, error){
	v, ok := ctx.Value(key).(Claims)
	if !ok {
		return Claims{}, errors.New("claim value missing from context")
	}
	return v, nil
}