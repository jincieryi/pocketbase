package core

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pocketbase/dbx"
)

func connectMysqlDB(mysqlDsn string) (*dbx.DB, error) {
	pragmas := "multiStatements=true"
	db, openErr := dbx.MustOpen("mysql", fmt.Sprintf("%s?%s", mysqlDsn, pragmas))
	return db, openErr
}
