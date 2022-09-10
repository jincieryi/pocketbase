package migrations_mysql

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/migrate"
)

var AppMigrations migrate.MigrationsList

func init() {
	AppMigrations.Register(
		func(db dbx.Builder) error {
			return nil
		},
		func(db dbx.Builder) error {
			return nil
		},
	)
}
