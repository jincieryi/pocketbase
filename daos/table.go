package daos

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/global"
)

// HasTable checks if a table with the provided name exists (case insensitive).
func (dao *Dao) HasTable(tableName string) bool {
	var exists bool
	if global.DB_TYPE == "mysql" {
		err := dao.DB().Select("count(*)").
			From("information_schema.TABLES").
			AndWhere(dbx.NewExp("LOWER([[TABLE_NAME]])=LOWER({:tableName})", dbx.Params{"tableName": tableName})).
			Limit(1).
			Row(&exists)

		return err == nil && exists
	} else {
		err := dao.DB().Select("count(*)").
			From("sqlite_schema").
			AndWhere(dbx.HashExp{"type": "table"}).
			AndWhere(dbx.NewExp("LOWER([[name]])=LOWER({:tableName})", dbx.Params{"tableName": tableName})).
			Limit(1).
			Row(&exists)

		return err == nil && exists
	}

}

// GetTableColumns returns all column names of a single table by its name.
func (dao *Dao) GetTableColumns(tableName string) ([]string, error) {
	columns := []string{}

	err := dao.DB().NewQuery("SELECT name FROM PRAGMA_TABLE_INFO({:tableName})").
		Bind(dbx.Params{"tableName": tableName}).
		Column(&columns)

	return columns, err
}

// DeleteTable drops the specified table.
func (dao *Dao) DeleteTable(tableName string) error {
	_, err := dao.DB().DropTable(tableName).Execute()

	return err
}
