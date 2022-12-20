package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/ethtweet/ethtweet/update"

	"github.com/ethtweet/ethtweet/appWeb/routes"
	"github.com/ethtweet/ethtweet/broadcastMsg"
	"github.com/ethtweet/ethtweet/config"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/logs"
	"github.com/ethtweet/ethtweet/migrates"
	"github.com/ethtweet/ethtweet/models"
	"github.com/ethtweet/ethtweet/p2pNet"
	"github.com/ethtweet/ethtweet/pRuntime"
	"github.com/ethtweet/ethtweet/tasks"

	"github.com/kataras/iris/v12"
	"github.com/sanbornm/go-selfupdate/selfupdate"
)

var updater = selfupdate.Updater{
	CurrentVersion: global.Version,                                             // Manually update the const, or set it using `go build -ldflags="-X main.VERSION=<newver>" -o ipfsTwitter src/ipfsTwitter/main.go`
	ApiURL:         "https://ipfstwitter.oss-cn-hongkong.aliyuncs.com/update/", // The server hosting `$CmdName/$GOOS-$ARCH.json` which contains the checksum for the binary
	BinURL:         "https://ipfstwitter.oss-cn-hongkong.aliyuncs.com/update/", // The server hosting the zip file containing the binary application which is a fallback for the patch method
	DiffURL:        "https://ipfstwitter.oss-cn-hongkong.aliyuncs.com/update/", // The server hosting the binary patch diff for incremental updates
	Dir:            "update/",                                                  // The directory created by the app when run which stores the cktime file
	CmdName:        "ethtweet",                                                 // The app name which is appended to the ApiURL to look for an update
	ForceCheck:     true,                                                       // For this example, always check for an update unless the version is "dev"
}

func init() {
	//需要在加载配置
	RunMysql()
	err := config.LoadConfig()
	if err != nil {
		logs.Fatal("reload config", err)
	}

	//注册需要编码传输的接口类型
	gob.Register(&broadcastMsg.TweetInfo{})
	gob.Register(&broadcastMsg.TweetInfoSync{})
	gob.Register(&broadcastMsg.UserInfo{})
	gob.Register(&p2pNet.HearBeat{})
}

func deleteOldFiles() {

	//window环境需要手动删除老版本文件
	filename := filepath.Base(os.Args[0])
	filename = "." + filename + ".old"
	_, err := os.Stat(filename)
	if err == nil {
		tmpFile := os.TempDir() + "EthTweet" + strconv.FormatInt(time.Now().UnixNano(), 10)
		err = os.Rename(filename, tmpFile)
		if err != nil {
			log.Println("delete old file err:", err)
			//window环境如果安装位置和盘符不一样，不能移动，先写在当前目录下的tmp目录
			os.Mkdir("./tmp/", os.FileMode(777))
			tmpFile = "./tmp/" + filename + ".old" + strconv.FormatInt(time.Now().UnixNano(), 10)
			os.Rename(filename, tmpFile)
		}
	}
}

func checkUpdate() {
	newVersion, _ := updater.UpdateAvailable()
	deleteOldFiles()
	if newVersion > global.Version {

		updater.Update()

		logs.Println("current version: ", global.Version)
		logs.Println("Update to version: ", newVersion)
		logs.Println("Ready to restart")
		time.Sleep(time.Second * 5) //更新前休眠5秒，避免重复冲突
		os.Exit(0)
	}
}

func checkUpdateTimer() {
	for {
		time.Sleep(time.Second * 600)
		logs.Println("checkUpdateTimer")
		checkUpdate()
		update.ChcckGithubVersion()
	}
}

func MigrateDb() error {
	db := global.GetDB()
	// 迁移 schema
	_ = migrates.Migrate()
	return db.AutoMigrate(
		&models.Tweets{},
		&models.User{},
		&models.Follow{},
		&models.UserAskSync{},
		&models.Tasks{},
	)
}

func BroadcastAll() {
	go func() {
		users := make([]*models.User, 0, 100)
		global.GetDB().Limit(100).Where("local_user = 1").Order("nonce desc").Offset(0).Find(&users)
		for _, u := range users {
			broadcastMsg.BroadcastNewestTweet(u)
			time.Sleep(time.Second * 10)
		}
	}()
	time.Sleep(time.Second * 10)
	users := make([]*models.User, 0, 100)
	global.GetDB().Limit(100).Where("local_nonce = nonce").Where("latest_cid != ''").Order("nonce desc").Offset(0).Find(&users)
	for _, u := range users {
		broadcastMsg.BroadcastNewestTweet(u)
		time.Sleep(time.Second * 10)
	}
}

func RunMysql() {
	if runtime.GOOS != "windows" {
		return
	}

	if err := global.CheckWindowsMysqld(); err != nil {
		logs.PrintDebugErr(err)
	}
}

func SavePeers() {

	conn := usr.Host.Network().Conns()
	//节点数过少，可能是网络中断等，暂停保存，避免覆盖
	if len(conn) < 50 {
		return
	}

	filePath := "Bootstrap.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	write := bufio.NewWriter(file)
	for _, c := range conn {
		//写入文件时，使用带缓存的 *Writer
		write.WriteString(fmt.Sprintf("%s/p2p/%s\n", c.RemoteMultiaddr().String(), c.RemotePeer().String()))
	}
	write.Flush()
}

var usr *p2pNet.UserNode

func main() {
	if config.Cfg.IsDaemon {
		if runtime.GOOS != "windows" {
			pRuntime.DaemonInit()
		} else {
			logs.PrintlnWarning("windows不支持守护进程模式")
		}
	}

	//不是debug模式下 自动管理进程
RE:
	proc, err := pRuntime.NewProc()
	if err != nil {
		logs.PrintlnWarning("up proc fail........")
	}
	//如果proc为nil表示当前进程已经是子进程了
	//不为空表示当前进程为主进程
	if proc != nil {
		go func() {
			pRuntime.HandleEndSignal(func() {
				if err := proc.Kill(); err != nil {
					logs.PrintErr(err)
				}
				logs.PrintlnSuccess("main proc exit....")

				os.Exit(0)
			})
		}()
		//等待子进程退出后 重启
		err = proc.Wait()
		if err != nil {
			logs.PrintlnWarning("proc wait err........")
		} else {
			goto RE
		}
		return
	} else {

		//CPU 性能分析
		if config.Cfg.Debug {
			go func() {
				log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
			}()
		}

		checkUpdate()
		update.ChcckGithubVersion()
		//子进程才执行更新检测
		go checkUpdateTimer()
	}

	go deleteOldFiles()
	go func() {
		err := global.ReloadIpfsGateway()
		if err != nil {
			logs.PrintErr(err)
		}
	}()

	fmt.Printf("EthTweet %s\n", global.Version)
	fmt.Printf("System version: %s\n", runtime.GOARCH+"/"+runtime.GOOS)
	fmt.Printf("Golang version: %s\n", runtime.Version())
	fmt.Printf("webui: http://localhost:%d/webui\n", config.Cfg.WebPort)

	// 迁移
	if err := MigrateDb(); err != nil {
		logs.Fatal(err)
	}

	usr = p2pNet.NewUserNode(config.Cfg.P2pPort, config.Cfg.UserKey, config.Cfg.KeyStore)
	err = usr.ConnectP2p()
	if err != nil {
		logs.Fatal("connect p2p err ", err)
	}
	logs.Println("your peer id ", usr.IDPretty())
	//设置当前用户
	err = models.SetCurrentUser(usr)

	if err != nil {
		logs.Fatal(err)
	}

	logs.PrintlnSuccess("your user id ", models.GetCurrentUser().Id)

	//询问最新用户资料
	logs.PrintlnInfo("ask user info ..............")
	go func() {
		err = broadcastMsg.SyncUserInfo(models.GetCurrentUser(), true)
		if err != nil {
			logs.PrintlnWarning("do ask err fail ", err)
		} else {
			logs.PrintlnSuccess("do ask is ok ")
		}
	}()

	//每
	ticker1 := time.NewTicker(300 * time.Second)
	// 一定要调用Stop()，回收资源
	defer ticker1.Stop()
	go func(t *time.Ticker) {
		for {
			// 每5秒中从chan t.C 中读取一次
			<-t.C
			SavePeers()
		}
	}(ticker1)

	go func() {
		tasks.RunTasks(global.GetGlobalCtx())
	}()
	go BroadcastAll()

	go func() {
		_ = broadcastMsg.SyncUserTweets(global.GetGlobalCtx())
	}()
	//可以监听系统上的ctrl+c信号 如果是守护进程模式 则在stop时触发appWeb
	appWeb := iris.New()

	go pRuntime.HandleEndSignal(func() {
		usr.Exit()
		models.ClearCurrentUser()
		err := appWeb.Shutdown(global.GetGlobalCtx())
		if err != nil {
			logs.PrintErr("appWeb shutdown err ", err)
		} else {
			logs.PrintlnSuccess("appWeb shutdown success")
		}
		global.CancelGlobalCtx()
		logs.PrintlnInfo("exit....")

		os.Exit(0)
	})

	//启用web服务器
	err = ListenWeb(appWeb)
	if err != nil {
		logs.Fatal(err)
	}
}

func ListenWeb(appWeb *iris.Application) (err error) {
	//注册api路由
	routes.RegisterApiRoutes(appWeb)

	logs.PrintlnInfo("Http API List:")
	for _, r := range appWeb.GetRoutes() {
		if r.Method != "OPTIONS" {
			logs.PrintlnInfo(fmt.Sprintf("[%s] http://127.0.0.1:%d%s", r.Method, config.Cfg.WebPort, r.Path))
		}
	}

	//监听http
	err = appWeb.Run(iris.Addr(fmt.Sprintf(":%d", config.Cfg.WebPort)), iris.WithConfiguration(iris.Configuration{
		TimeFormat: global.DateTimeFormatStr,
	}))
	if err != nil {
		return err
	}
	return nil
}
