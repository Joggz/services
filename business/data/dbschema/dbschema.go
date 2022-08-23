// Package dbschema contains the database schema, migrations and seeding data.
package dbschema

import (
	"context"
	_ "embed" // Calls init function.
	"fmt"

	"github.com/Joggz/services/business/sys/database"
	"github.com/ardanlabs/darwin"
	"github.com/jmoiron/sqlx"
)

var (
	//go:embed sql/schema.sql
	schemaDoc string

	//go:embed sql/seed.sql
	seedDoc string

	//go:embed sql/delete.sql
	deleteDoc string
)



func Migrate(ctx context.Context, db *sqlx.DB) error {
	if	err :=	database.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)	
	}

	driver, err := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})

	if err != nil {
		return fmt.Errorf("conrusting darwin driver %w", err)
	}

	d := darwin.New(driver, darwin.ParseMigrations(schemaDoc))

	return d.Migrate()

}

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func Seed(ctx context.Context, db *sqlx.DB) error {
	if	err :=	database.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)	
	}

	tx, err := db.Begin()
	if err!=nil {
		return err
	}

	if _, err := tx.Exec(seedDoc); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
	}


	return tx.Commit()
}