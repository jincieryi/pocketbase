package migrations_mysql

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

const (
	// CollectionsExtendCollectionName is the name of the system user collectionsExtend collection.
	CollectionsExtendCollectionName = "collectionsExtend"
)

//创建 datasources connection
func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		// inserts the system collectionsExtend collection
		// -----------------------------------------------------------

		collection := &models.Collection{
			Name:   CollectionsExtendCollectionName,
			System: true,
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Id:       "pbfieldcid",
					Name:     "cid",
					Unique:   true,
					Required: true,
					Type:     schema.FieldTypeText,
					Options:  &schema.TextOptions{},
				},
				&schema.SchemaField{
					Id:       "pbfielddid",
					Name:     "did",
					Required: true,
					Type:     schema.FieldTypeText,
					Options:  &schema.TextOptions{},
				},
				&schema.SchemaField{
					Id:       "pbfieldrawsql",
					Name:     "rawSql",
					Required: false,
					Type:     schema.FieldTypeText,
					Options:  &schema.TextOptions{},
				},
				&schema.SchemaField{
					Id:       "pbfieldidname",
					Name:     "idName",
					Required: false,
					Type:     schema.FieldTypeText,
					Options:  &schema.TextOptions{},
				},
				&schema.SchemaField{
					Id:       "pbfieldextend",
					Name:     "extend",
					Required: false,
					Type:     schema.FieldTypeJson,
					Options:  &schema.JsonOptions{},
				},
			),
		}
		collection.Id = "systemcollectionsextend0"
		collection.MarkAsNew()
		collection.System = true

		return daos.New(db).SaveCollection(collection)

	}, func(db dbx.Builder) error {
		_, err := db.DropTable(CollectionsExtendCollectionName).Execute()
		return err
	})
}
