package config

import (
	"encoding/json"
	"flag"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/logs"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

var webPort = flag.Int("web_port", 8080, "web listen port")
var p2pPort = flag.Int("port", 20000, "libp2p listen port")
var userKey = flag.String("user_key", "userKey77", "user keystore")
var ipfsApi = flag.String("ipfs_api", "https://cdn.ipfsscan.io", "ipfs api")
var keyStore = flag.String("key_store", "./keyStore", "keystore dir")
var userData = flag.String("user_data", global.SqliteDatabaseDefaultDir, "data dir")
var debug = flag.Bool("debug", false, "is print logs")
var isLog = flag.Bool("is_print", true, "is print debug logs")
var isCheckApiLocal = flag.Bool("check_api_local", true, "is check api local")
var isDaemon = flag.Bool("isDaemon", false, "启用守护进程模式 不支持win")
var configPath = flag.String("config", "", "配置文件路径")

var Cfg *Conf

type Conf struct {
	IsDaemon      bool
	Logs          *LogsCfg
	Debug         bool
	CheckApiLocal bool
	UserKey       string
	KeyStore      string
	WebPort       int
	P2pPort       int
	IpfsApi       string
}

type LogsCfg struct {
	IsPrint     bool
	LogFilePath string
	LogFile     *os.File
}

func init() {
	flag.Parse()
}

func LoadConfig() error {
	logs.PrintlnWarning(*configPath)
	return LoadConfigByPath(*configPath)
}

func LoadConfigByPath(p string) error {
	defer func() {
		logs.PrintlnSuccess("load config success!")
		logs.IsPrintLog = Cfg.Logs.IsPrint
		logs.IsDebugPrint = Cfg.Debug
	}()
	_, err := os.Stat("tweet.yaml")
	if err == nil && p == "" {
		p = "./tweet.yaml"
	}
	if p == "" {
		logs.PrintlnInfo("load config from cli args...")
		Cfg = &Conf{
			IsDaemon: *isDaemon,
			Logs: &LogsCfg{
				IsPrint:     *isLog,
				LogFilePath: "",
				LogFile:     nil,
			},
			Debug:         *debug,
			CheckApiLocal: *isCheckApiLocal,
			UserKey:       *userKey,
			IpfsApi:       *ipfsApi,
			KeyStore:      *keyStore,
			WebPort:       *webPort,
			P2pPort:       *p2pPort,
		}
		global.DbDrive = global.DBDriveSqlite
		return global.InitSqliteDatabase(*userData, global.SqliteDatabaseName)
	} else {
		logs.PrintlnInfo("load config from config file...", *configPath)
	}
	Cfg = &Conf{}
	paths, fileName := filepath.Split(p)
	viper.SetConfigName(fileName)
	viper.AddConfigPath(paths)
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	drive := viper.GetString("db.drive")
	if drive == global.DBDriveMysql {
		b := make([]byte, 0)
		b, err := json.Marshal(viper.GetStringMap("db.mysql"))
		if err != nil {
			return err
		}
		dbConf := gjson.ParseBytes(b)
		maxIdleCounts, _ := strconv.Atoi(dbConf.Get("max_idle_counts").String())
		charset := dbConf.Get("charset").String()
		if charset == "" {
			charset = "utf8"
		}
		db, err := global.NewDatabaseMysql(dbConf.Get("host").String(), dbConf.Get("port").String(), dbConf.Get("database").String(), charset, dbConf.Get("username").String(), dbConf.Get("password").String(), maxIdleCounts, int(dbConf.Get("max_open_counts").Int()))
		if err != nil {
			return err
		}
		global.DbDrive = global.DBDriveMysql
		global.SetMysqlDB(db)
	} else {
		global.DbDrive = global.DBDriveSqlite
		dir := viper.GetString("db.sqlite.dir")
		if dir == "" {
			dir = *userData
		}
		if runtime.GOOS == "android" {
			dir = dir + dir
		}
		dbName := viper.GetString("db.sqlite.dbName")
		err = global.InitSqliteDatabase(dir, dbName)
		if err != nil {
			return err
		}
	}

	//载入日志配置
	logCfg := &LogsCfg{}
	logCfg.IsPrint = viper.GetBool("log.isPrint")
	logCfg.LogFilePath = viper.GetString("log.logFilePath")
	if logCfg.LogFilePath != "" {
		logCfg.LogFile, err = os.OpenFile(logCfg.LogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		os.Stdout = logCfg.LogFile
		os.Stderr = logCfg.LogFile
		go func() {
			t := time.NewTicker(36 * time.Hour)
			defer t.Stop()
			log.Println("创建自动清除日志文件内容")
			for {
				select {
				case <-t.C:
					err = logCfg.LogFile.Truncate(0)
					if err != nil {
						log.Println("清空文件失败")
					}
					_, _ = logCfg.LogFile.Seek(0, 0)

				}
			}
		}()
	}
	Cfg.Logs = logCfg
	Cfg.Debug = viper.GetBool("debug")
	Cfg.CheckApiLocal = viper.GetBool("check_api_local")
	Cfg.IsDaemon = viper.GetBool("daemon")
	Cfg.UserKey = viper.GetString("key.user_key")
	Cfg.IpfsApi = viper.GetString("ipfs_api")
	Cfg.KeyStore = viper.GetString("key.key_store")
	Cfg.P2pPort = viper.GetInt("port.p2p_port")
	Cfg.WebPort = viper.GetInt("port.web_port")
	if Cfg.P2pPort == 0 {
		Cfg.P2pPort = *p2pPort
	}
	if Cfg.WebPort == 0 {
		Cfg.WebPort = *webPort
	}
	if Cfg.UserKey == "" {
		Cfg.UserKey = *userKey
	}
	if Cfg.KeyStore == "" {
		Cfg.KeyStore = *userData
	}
	if Cfg.IpfsApi == "" {
		Cfg.IpfsApi = "https://cdn.ipfsscan.io"
	}
	return nil
}
