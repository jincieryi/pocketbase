package core

import (
	"github.com/go-sql-driver/mysql"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/tools/store"
)

type MountDBProvider struct {
	masterDao *daos.Dao
	cache     *store.Store[*MountDB]
}

type MountDB struct {
	id         string
	db         *dbx.DB
	schemaName string
}

//初始化MountDBProvider
func NewMountDBProvider(masterDao *daos.Dao) *MountDBProvider {
	return &MountDBProvider{masterDao: masterDao, cache: store.New[*MountDB](nil)}
}

//构造MountDB
func NewMountDB(id string, dsn string) (*MountDB, error) {
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	db, err := connectMysqlDB(dsn)
	if err != nil {
		return nil, err
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

	db, err := NewMountDB(id, datasource.Data()["dsn"].(string))
	if err != nil {
		return nil, err
	}

	provider.Put(id, db)

	return db, nil
}

func (provider *MountDBProvider) Remove(id string) {
	provider.cache.Remove(id)
}
