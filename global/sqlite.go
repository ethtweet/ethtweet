package global

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

const SqliteDatabaseName = "EthTweet.0.2.db"
const SqliteDatabaseDefaultDir = "databases"

var sqliteDb *SqliteDB

type SqliteDB struct {
	dbName string
	dir    string
	*gorm.DB
}

func (sdb *SqliteDB) GetDsn() string {
	return sdb.dir + "/" + sdb.dbName + "?_journal_mode=WAL"
}

func GetSqliteDB() *SqliteDB {
	return sqliteDb
}

func InitSqliteDatabase(dir, name string) error {
	var err error
	if name == "" {
		name = SqliteDatabaseName
	}
	sqliteDb = &SqliteDB{
		dbName: name,
		dir:    dir,
	}
	if !IsDir(sqliteDb.dir) {
		err = os.Mkdir(sqliteDb.dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	sqliteDb.DB, err = gorm.Open(sqlite.Open(sqliteDb.GetDsn()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	sqliteDb.DB.Exec("PRAGMA SQLITE_THREADSAFE=2")
	sqliteDb.DB.Exec("PRAGMA foreign_keys = ON")
	sqliteDb.DB.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return err
	}
	return nil
}
