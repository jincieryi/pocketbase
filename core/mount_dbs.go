package core

import (
	"context"
	"database/sql"
	"github.com/fatih/color"
	"github.com/go-sql-driver/mysql"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/tools/store"
	"time"
)

type MountDBProvider struct {
	masterDao *daos.Dao
	cache     *store.Store[*MountDB]
	isDebug   bool
}

type MountDB struct {
	id         string
	db         *dbx.DB
	schemaName string
}

//初始化MountDBProvider
func NewMountDBProvider(masterDao *daos.Dao, isDebug bool) *MountDBProvider {
	return &MountDBProvider{masterDao: masterDao, cache: store.New[*MountDB](nil), isDebug: isDebug}
}

//构造MountDB
func newMountDB(id string, dsn string, isDebug bool) (*MountDB, error) {
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	db, err := connectMysqlDB(dsn)
	if err != nil {
		return nil, err
	}

	if isDebug {
		db.QueryLogFunc = func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
			color.HiBlack("[MountDB][%.2fms] %v\n", float64(t.Milliseconds()), sql)

			if err != nil {
				color.HiRed(err.Error())
			}
		}
	}

	return &MountDB{
		id: id, db: db, schemaName: config.DBName,
	}, err
}

func (this *MountDB) GetDB() *dbx.DB {
	return this.db
}

func (provider *MountDBProvider) Put(id string, mdb *MountDB) {
	provider.cache.Set(id, mdb)
}

func (provider *MountDBProvider) Get(id string) (*MountDB, error) {

	if mountdb := provider.cache.Get(id); mountdb != nil {
		return mountdb, nil
	}

	collection, _ := provider.masterDao.FindCollectionByNameOrId("datasources")
	datasource, err := provider.masterDao.FindRecordById(collection, id, nil)
	if err != nil {
		return nil, err
	}

	db, err := newMountDB(id, datasource.Data()["dsn"].(string), provider.isDebug)
	if err != nil {
		return nil, err
	}

	provider.Put(id, db)

	return db, nil
}

func (provider *MountDBProvider) Remove(id string) {
	provider.cache.Remove(id)
}
