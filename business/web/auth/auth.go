package auth

import (
	"crypto/rsa"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

// Auth is used to authenticate clients. It can generate a token for a
// set of user claims and recreate the claims by parsing the token.
type Auth struct {
	activeKeyID string	
	KeyLookUp KeyLookUp
	method jwt.SigningMethod
	KeyFunc func (t *jwt.Token) (any, error)
	parser *jwt.Parser
}

// KeyLookup declares a method set of behavior for looking up
// private and public keys for JWT use.
type KeyLookUp interface {
	PrivateKey(kid string) (*rsa.PrivateKey, error)
	Publickey(kid string) (*rsa.PublicKey, error)
}

// New creates an Auth to support authentication/authorization.
func New(activeKeyID string,  KeyLookUp KeyLookUp) (*Auth, error) {
	// The activeKID represents the private key used to signed new tokens.
	_, err := KeyLookUp.PrivateKey(activeKeyID)
	if err != nil {
		return nil, errors.New("active KID does not exist in store")
	}
	method := jwt.GetSigningMethod("RS256");
	if method == nil {
		return nil, errors.New("configuring algorithm RS256")
	}
	keyFunc := func(t *jwt.Token) (any, error) {
		kid, ok := t.Header["kid"]
		if !ok {
			return nil, errors.New("missing key id (kid) in token header")
		}

		KidID, ok := kid.(string)
		if !ok {
			return nil, errors.New("user token key id (kid) must be string")
		}
		return KeyLookUp.Publickey(KidID)
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))

	
	a := Auth{
		activeKeyID: activeKeyID,
	KeyLookUp: KeyLookUp,
	method: method,
	KeyFunc: keyFunc,
	parser: parser,
	}

	return &a, nil
}