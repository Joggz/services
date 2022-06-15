// Package dbschema contains the database schema, migrations and seeding data.
package dbschema

import (
	_ "embed" // Calls init function.

	"github.com/ardanlabs/darwin"
)

var (
	//go:embed sql/schema.sql
	schemaDoc string

	//go:embed sql/seed.sql
	seedDoc string

	//go:embed sql/delete.sql
	deleteDoc string
)

func D ()  {
	darwin.New()
}
