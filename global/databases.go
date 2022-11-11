package global

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	DBDriveSqlite = "sqlite"
	DBDriveMysql  = "mysql"
)

var DbDrive = ""

func init() {
	if DbDrive == "" {
		DbDrive = DBDriveSqlite
	}
}

func GetDB() *gorm.DB {
	if DbDrive == DBDriveSqlite {
		return sqliteDb.DB
	}
	return mysqlDb.DB
}

func LockForUpdate(db *gorm.DB) *gorm.DB {
	//sqlite不加锁
	/*if db.Name() == DBDriveSqlite {
		log.Println("sqlite not lock.......")
		return db
	}*/
	return db.Clauses(clause.Locking{Strength: "UPDATE"})
}
