// Package dbtest contains supporting code for running tests that hit the DB.
package dbtest

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"os"

	// "os/user"
	"testing"
	"time"

	"github.com/Joggz/services/business/data/dbschema"
	"github.com/Joggz/services/business/data/store/user"
	"github.com/Joggz/services/business/sys/database"
	"github.com/Joggz/services/business/web/auth"
	"github.com/Joggz/services/foundation/docker"
	"github.com/Joggz/services/foundation/keystore"
	"github.com/Joggz/services/foundation/logger"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

type DBContainer struct {
	Image string
	Port string
	Args []string
}



// NewUnit creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty. It returns
// the database to use as well as a function to call at the end of the test.
func NewUnit(t *testing.T, dbc DBContainer)(*zap.SugaredLogger, *sqlx.DB, func()){
		r, w, _ := os.Pipe()
		old  := os.Stdout
		os.Stdout = w

		c, _ := docker.StartContainer(dbc.Image, dbc.Port, dbc.Args...)

		db, err := database.Open(database.Config{
			User:         "postgres",
			Password:     "password",
			Host:         c.Host,
			Name:         "postgres",
			DisableTLS:   true,
		})

		if err !=nil {
			t.Fatalf("opening database connection: %v", err)
		}

		// t.Log("waiting for database to be ready", db )

	ctx, cancel :=	context.WithTimeout(context.Background(), 30*time.Second)
	t.Log("waiting for database to be ready", db )
	
	defer cancel()
	t.Log("database should be ready", db )
	if err :=	dbschema.Migrate(ctx, db); err != nil {
		docker.DumpContainerLogs(t, c.ID)
		docker.StopContainer(c.ID)
		t.Fatalf("migration error: %v", err)
	}

	if err := dbschema.Seed(ctx, db); err != nil {
		docker.DumpContainerLogs(t, c.ID)
		docker.StopContainer(c.ID)
		t.Fatalf("seeding error: %v", err)
	}

	log, err := logger.New("TEST")
	if err != nil {
		t.Fatalf("logger error: %v", err)
	}
	teardown := func(){
		t.Helper()
		db.Close()
		docker.StopContainer(c.ID)

		log.Sync()
		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = old

		fmt.Println("******************** LOGS ********************")
		fmt.Print(buf.String())
		fmt.Println("******************** LOGS ********************")

	}
	return log, db, teardown
}

// Test owns state for running and shutting down tests.
type Test struct {
	DB       *sqlx.DB
	Log      *zap.SugaredLogger
	Auth     *auth.Auth
	Teardown func()

	t *testing.T
}


// NewIntegration creates a database, seeds it, constructs an authenticator.
func NewIntegration(t *testing.T, dbc DBContainer) *Test {
	log, db, teardown := NewUnit(t, dbc)

		// Create RSA keys to enable authentication in our service.
		keyID := "4754d86b-7a6d-4df5-9c65-224741361492"
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}
	
		// Build an authenticator using this private key and id for the key store.
		auth, err := auth.New(keyID, keystore.NewMap(map[string]*rsa.PrivateKey{keyID: privateKey}))
		if err != nil {
			t.Fatal(err)
		}
	
	test := Test{
		DB:   db,
		Log:  log,
		Auth: auth,
		Teardown: teardown,
		t: t,
	}
	return &test
}

// Token generates an authenticated token for a user.
func (test *Test) Token(email, pass string) string {
	test.t.Log("Generating token for test ...")

	store := user.NewStore(test.Log, test.DB)
	claim, err :=	store.Authenticate(context.Background(), time.Now(), email, pass)
	if err != nil {
	test.t.Fatal(err)
	}

	token, err := test.Auth.GenerateToken(claim)
	if err != nil {
		test.t.Fatal(err)
	}

	return token
}


func StringPointer(s string) *string {
	return &s
}