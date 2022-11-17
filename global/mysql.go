package global

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var mysqlDb *MysqlDb

type MysqlDb struct {
	Host          string
	Port          string
	Database      string
	Charset       string
	UserName      string
	Password      string
	MaxIdleCounts int
	*gorm.DB
}

func GetMysqlDB() *MysqlDb {
	return mysqlDb
}

func SetMysqlDB(db *MysqlDb) {
	mysqlDb = db
}

func NewDatabaseMysql(host, port, database, charset, username, password string, maxIdleCounts int, maxOpenCounts int) (*MysqlDb, error) {

	// 初始化数据库
	if runtime.GOOS == "windows" {
		_, err := os.Stat("mysql/data/ibdata1")
		_, err2 := os.Stat("mysql/data/" + database + "/tweets.ibd")
		if err == nil && err2 != nil {
			db, err := sql.Open("mysql", username+":"+password+"@("+host+":"+port+")/mysql?charset="+charset+"&parseTime=True&loc=Local")
			if err != nil {
				fmt.Println("mysql Open", err)
				//return nil, err
			}
			defer db.Close()

			_, err = db.Exec("CREATE DATABASE " + database + ";")
			if err != nil {
				fmt.Println("mysql CREATE ", err)
			}
		}
	}

	db, err := gorm.Open(mysql.Open(username+":"+password+"@("+host+":"+port+")/"+database+"?charset="+charset+"&parseTime=True&loc=Local"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold:             time.Second,  // 慢 SQL 阈值
				LogLevel:                  logger.Error, // 日志级别
				IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,
			},
		),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if maxIdleCounts > 0 {
		sqlDB.SetMaxIdleConns(maxIdleCounts)
	}
	if maxOpenCounts > 0 {
		sqlDB.SetMaxOpenConns(maxOpenCounts)
	}
	sqlDB.SetConnMaxLifetime(1 * time.Minute)
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}
	return &MysqlDb{
		Host:          host,
		Port:          port,
		Database:      database,
		Charset:       charset,
		UserName:      username,
		Password:      password,
		MaxIdleCounts: maxIdleCounts,
		DB:            db,
	}, nil
}
