package migrations_mysql

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

const (
	// DatasourceCollectionName is the name of the system user datasources collection.
	DatasourceCollectionName = "datasources"
)

//创建 datasources connection
func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		// inserts the system datasources collection
		// -----------------------------------------------------------

		collection := &models.Collection{
			Name:   DatasourceCollectionName,
			System: true,
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Id:       "pbfieldname",
					Name:     "name",
					Unique:   true,
					Required: true,
					Type:     schema.FieldTypeText,
					Options:  &schema.TextOptions{},
				},
				&schema.SchemaField{
					Id:       "pbfielddsn",
					Name:     "dsn",
					Required: true,
					Type:     schema.FieldTypeText,
					Options:  &schema.TextOptions{},
				},
				&schema.SchemaField{
					Id:       "pbfieldtype",
					Name:     "type",
					Required: true,
					Type:     schema.FieldTypeSelect,
					Options: &schema.SelectOptions{
						MaxSelect: 1,
						Values:    []string{"mysql"},
					},
				},
			),
		}
		collection.Id = "systemdatasources0"
		collection.MarkAsNew()
		collection.System = true

		return daos.New(db).SaveCollection(collection)

	}, func(db dbx.Builder) error {
		_, err := db.DropTable(DatasourceCollectionName).Execute()
		return err
	})
}
