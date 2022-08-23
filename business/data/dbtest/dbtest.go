// Package dbtest contains supporting code for running tests that hit the DB.
package dbtest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/Joggz/services/business/data/dbschema"
	"github.com/Joggz/services/business/sys/database"
	"github.com/Joggz/services/foundation/docker"
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

func StringPointer(s string) *string {
	return &s
}