package global

import (
	"gorm.io/gorm"
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
