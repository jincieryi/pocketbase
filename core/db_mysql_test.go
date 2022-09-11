package core

import (
	"context"
	"database/sql"
	"github.com/fatih/color"
	"testing"
	"time"
)

func TestConnectMysqlDB(t *testing.T) {
	_, err := connectMysqlDB("root:root@tcp(1.117.39.176:3306)/pb_test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateTable(t *testing.T) {
	db, err := connectMysqlDB("root:root@tcp(1.117.39.176:3306)/pb_test")
	if err != nil {
		t.Fatal(err)
	}

	db.ExecLogFunc = func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		color.HiBlack("[%.2fms] %v\n", float64(t.Milliseconds()), sql)
	}

	result, err := db.NewQuery(`

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
				[[email]]                  VARCHAR(100) UNIQUE NOT NULL,
				[[tokenKey]]               VARCHAR(100) UNIQUE NOT NULL,
				[[passwordHash]]           VARCHAR(100) NOT NULL,
				[[lastResetSentAt]]        VARCHAR(100) DEFAULT "" NOT NULL,
				[[lastVerificationSentAt]] VARCHAR(100) DEFAULT "" NOT NULL,
				[[created]]                VARCHAR(100) DEFAULT "" NOT NULL,
				[[updated]]                VARCHAR(100) DEFAULT "" NOT NULL
			);

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

`).Execute()

	if err != nil {
		t.Fatal(err)
	}

	println(result)

}
