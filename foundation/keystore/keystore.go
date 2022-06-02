// Package keystore implements the auth.KeyLookup interface. This implements
// an in-memory keystore for JWT support.

package keystore

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v4"
)

// KeyStore represents an in memory store implementation of the
// KeyLookup interface for use with the auth package.
type KeyStore struct {
	mu sync.RWMutex
	store map[string]*rsa.PrivateKey
}

// New constructs an empty KeyStore ready for use.
func New() *KeyStore{
	return &KeyStore{
		store: make(map[string]*rsa.PrivateKey),
	}
}

// NewMap constructs a KeyStore with an initial set of keys.
func NewMap(store map[string]*rsa.PrivateKey) *KeyStore {
	return &KeyStore{
		store: store,
	}
}

func NewFS(fys fs.FS) (*KeyStore, error) {
	ks := New()

	fn := func(fileName string, dirEntry fs.DirEntry, err error) error{
		if err != nil {
			return fmt.Errorf("walkdir failure: %w", err)
		}

		if dirEntry.IsDir(){
			return nil
		}

		file, err := fys.Open(fileName)
		if err != nil {
			return fmt.Errorf("opening key file: %w", err)
		}
		if path.Ext(fileName) != ".pem" {
			return nil
		}
		defer file.Close()

		// limit PEM file size to 1 megabyte. This should be reasonable for
		// almost any PEM file and prevents shenanigans like linking the file
		// to /dev/random or something like that.

		privatePem, err := io.ReadAll(io.LimitReader(file, 1024*1024))
		if err != nil {
			return fmt.Errorf("reading auth private key: %w", err)
		}
		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePem)
		if err != nil {
			return fmt.Errorf("parsing auth private key: %w", err)
		}
		
		ks.store[strings.TrimSuffix(dirEntry.Name(), ".pem")] = privateKey
		
		return nil
	}


	 err := fs.WalkDir(fys, ".", fn)
	 if err !=  nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	 }
	return ks, nil
}

// Add adds a private key and combination kid to the store.
func (ks *KeyStore) Add(privateKey *rsa.PrivateKey, kid string) {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	ks.store[kid] = privateKey
}

// Remove removes a private key and combination kid to the store.
func (ks *KeyStore) Remove(kid string) {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	delete(ks.store, kid)
}

// PrivateKey searches the key store for a given kid and returns
// the private key.
func (ks *KeyStore) PrivateKey(kid string) (*rsa.PrivateKey, error) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	privateKey, found := ks.store[kid]
	if !found {
		return nil, errors.New("kid lookup failed")
	}
	return privateKey, nil
}

// PublicKey searches the key store for a given kid and returns
// the public key.
func (ks *KeyStore) PublicKey(kid string) (*rsa.PublicKey, error) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	privateKey, found := ks.store[kid]
	if !found {
		return nil, errors.New("kid lookup failed")
	}
	return &privateKey.PublicKey, nil
}
