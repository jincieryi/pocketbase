package migrations_mysql

import (
	"fmt"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/migrate"
)

var AppMigrations migrate.MigrationsList

func init() {
	AppMigrations.Register(
		func(db dbx.Builder) error {
			_, tablesErr := db.NewQuery(`

			CREATE TABLE IF NOT EXISTS {{_admins}} (
				[[id]]              VARCHAR(100) PRIMARY KEY,
				[[avatar]]          INTEGER DEFAULT 0 NOT NULL,
				[[email]]           VARCHAR(100) UNIQUE NOT NULL,
				[[tokenKey]]        VARCHAR(100) UNIQUE NOT NULL,
				[[passwordHash]]    VARCHAR(100) NOT NULL,
				[[lastResetSentAt]] VARCHAR(100) DEFAULT "" NOT NULL,
				[[created]]         VARCHAR(100) DEFAULT "" NOT NULL,
				[[updated]]         VARCHAR(100) DEFAULT "" NOT NULL
			);

			CREATE TABLE IF NOT EXISTS {{_users}} (
				[[id]]                     VARCHAR(100) PRIMARY KEY,
				[[verified]]               TINYINT DEFAULT 0 NOT NULL,
				[[email]]                  VARCHAR(100) DEFAULT "" NOT NULL,
				[[tokenKey]]               VARCHAR(100) NOT NULL,
				[[passwordHash]]           VARCHAR(100) NOT NULL,
				[[lastResetSentAt]]        VARCHAR(100) DEFAULT "" NOT NULL,
				[[lastVerificationSentAt]] VARCHAR(100) DEFAULT "" NOT NULL,
				[[created]]                VARCHAR(100) DEFAULT "" NOT NULL,
				[[updated]]                VARCHAR(100) DEFAULT "" NOT NULL
			);

			CREATE UNIQUE INDEX _users_email_idx ON {{_users}} ([[email]]);
			CREATE UNIQUE INDEX _users_tokenKey_idx ON {{_users}} ([[tokenKey]]);

			CREATE TABLE IF NOT EXISTS {{_collections}} (
				[[id]]         VARCHAR(100) PRIMARY KEY,
				[[system]]     TINYINT DEFAULT 0 NOT NULL,
				[[name]]       VARCHAR(100) UNIQUE NOT NULL,
				[[schema]]     JSON  NOT NULL,
				[[listRule]]   VARCHAR(100) DEFAULT NULL,
				[[viewRule]]   VARCHAR(100) DEFAULT NULL,
				[[createRule]] VARCHAR(100) DEFAULT NULL,
				[[updateRule]] VARCHAR(100) DEFAULT NULL,
				[[deleteRule]] VARCHAR(100) DEFAULT NULL,
				[[created]]    VARCHAR(100) DEFAULT "" NOT NULL,
				[[updated]]    VARCHAR(100) DEFAULT "" NOT NULL
			);

			CREATE TABLE IF NOT EXISTS {{_params}} (
				[[id]]      VARCHAR(100) PRIMARY KEY,
				[[key]]     VARCHAR(100) UNIQUE NOT NULL,
				[[value]]   JSON DEFAULT NULL,
				[[created]] VARCHAR(100) DEFAULT "" NOT NULL,
				[[updated]] VARCHAR(100) DEFAULT "" NOT NULL
			);

			CREATE TABLE IF NOT EXISTS {{_externalAuths}} (
				[[id]]         VARCHAR(100) PRIMARY KEY,
				[[userId]]     VARCHAR(100) NOT NULL,
				[[provider]]   VARCHAR(100) NOT NULL,
				[[providerId]] VARCHAR(100) NOT NULL,
				[[created]]    VARCHAR(100) DEFAULT "" NOT NULL,
				[[updated]]    VARCHAR(100) DEFAULT "" NOT NULL,
			    constraint FOREIGN KEY ([[userId]]) REFERENCES {{_users}} ([[id]]) ON UPDATE CASCADE ON DELETE CASCADE
			);

			CREATE UNIQUE INDEX _externalAuths_userId_provider_idx on {{_externalAuths}} ([[userId]], [[provider]]);
			CREATE UNIQUE INDEX _externalAuths_provider_providerId_idx on {{_externalAuths}} ([[provider]], [[providerId]]);

			`).Execute()

			if tablesErr != nil {
				return tablesErr
			}
			// inserts the system profiles collection
			// -----------------------------------------------------------
			profileOwnerRule := fmt.Sprintf("%s = @request.user.id", models.ProfileCollectionUserFieldName)
			collection := &models.Collection{
				Name:       models.ProfileCollectionName,
				System:     true,
				CreateRule: &profileOwnerRule,
				ListRule:   &profileOwnerRule,
				ViewRule:   &profileOwnerRule,
				UpdateRule: &profileOwnerRule,
				Schema: schema.NewSchema(
					&schema.SchemaField{
						Id:       "pbfielduser",
						Name:     models.ProfileCollectionUserFieldName,
						Type:     schema.FieldTypeUser,
						Unique:   true,
						Required: true,
						System:   true,
						Options: &schema.UserOptions{
							MaxSelect:     1,
							CascadeDelete: true,
						},
					},
					&schema.SchemaField{
						Id:      "pbfieldname",
						Name:    "name",
						Type:    schema.FieldTypeText,
						Options: &schema.TextOptions{},
					},
					&schema.SchemaField{
						Id:   "pbfieldavatar",
						Name: "avatar",
						Type: schema.FieldTypeFile,
						Options: &schema.FileOptions{
							MaxSelect: 1,
							MaxSize:   5242880,
							MimeTypes: []string{
								"image/jpg",
								"image/jpeg",
								"image/png",
								"image/svg+xml",
								"image/gif",
							},
						},
					},
				),
			}
			collection.Id = "systemprofiles0"
			collection.MarkAsNew()

			return daos.New(db).SaveCollection(collection)
		},
		func(db dbx.Builder) error {
			tables := []string{
				"_params",
				"_collections",
				"_users",
				"_admins",
				"_externalAuths",
				models.ProfileCollectionName,
			}

			for _, name := range tables {
				if _, err := db.DropTable(name).Execute(); err != nil {
					return err
				}
			}

			return nil
		})
}
