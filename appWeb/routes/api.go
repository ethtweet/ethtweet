package routes

import (
	"fmt"
	"github.com/ethtweet/ethtweet/appWeb/controller"
	"github.com/ethtweet/ethtweet/appWeb/middleware"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/models"
	"github.com/ethtweet/ethtweet/update"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func RegisterApiRoutes(app *iris.Application) {
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	var templatesDir = "./templates"
	if runtime.GOOS == "android" {
		ex, err := os.Executable()
		if err != nil {
			fmt.Println(err)
		}
		exPath := filepath.Dir(ex)

		templatesDir = exPath + "/templates"
		err = update.Unzip(exPath+"/templates.zip", templatesDir)
		if err != nil {
			fmt.Println("templatesDir Unzip:" + err.Error())
		}
	}

	tmpl := iris.HTML(templatesDir, ".html")
	tmpl.AddFunc("Split", func(s string) []string {
		return strings.Split(s, ",")
	})
	tmpl.AddFunc("FormatTime", func(s int64) string {
		return time.Unix(s, 0).Format(time.DateTime)
	})

	tmpl.AddFunc("image", func(s string) string {
		return "https://image.baidu.com/search/down?url=" + s
	})

	app.RegisterView(tmpl)

	app.Get("/webui", func(ctx iris.Context) {
		ctx.Redirect("https://ipfs.io/ipns/share.ethtweet.io")
	})

	app.HandleDir("static", templatesDir)

	app.Get("/", func(ctx iris.Context) {
		//ctx.Writef("Hello from method: %s and path: %s id:%s", ctx.Method(), ctx.Path(), ctx.Params().Get("id"))
		// 绑定： {{.message}}　为　"Hello world!"
		// 渲染模板文件： ./views/hello.html
		pager := global.NewPager(ctx)
		tws := make([]*models.Tweets, 0, pager.Limit)
		global.GetDB().Limit(pager.Limit).Offset(pager.Offset).Preload("UserInfo").Order("created_at desc").Find(&tws)
		ctx.ViewData("tweets", tws)
		ctx.View("index.html")
	})

	app.Get("/user/{id:string}", func(ctx iris.Context) {
		//ctx.Writef("Hello from method: %s and path: %s id:%s", ctx.Method(), ctx.Path(), ctx.Params().Get("id"))
		// 绑定： {{.message}}　为　"Hello world!"
		// 渲染模板文件： ./views/hello.html
		id := ctx.Params().Get("id")
		user := &models.User{}
		global.GetDB().Model(user).Where("id = ?", id).Find(&user)
		pager := global.NewPager(ctx)
		tws := make([]*models.Tweets, 0, pager.Limit)
		global.GetDB().Where("user_id", id).Limit(pager.Limit).
			Preload("OriginTw").
			Preload("UserInfo").
			Order("nonce desc").
			Offset(pager.Offset).Find(&tws)
		ctx.ViewData("user", user)
		ctx.ViewData("tweets", tws)
		ctx.ViewData("id", id)
		ctx.View("user-page.html")
	})

	v0 := app.Party("/api/v0", crs).AllowMethods(iris.MethodOptions)
	{
		mvc.Configure(v0.Party("/"), func(application *mvc.Application) {
			application.Handle(&controller.SiteController{})
		})

		mvc.New(v0.Party("/", middleware.ApiLocalAuth)).Register(models.GetCurrentUser()).Configure(func(application *mvc.Application) {
			application.Party("/tweet").Handle(&controller.TweetsController{})
			application.Party("/user").Handle(&controller.UserController{})
			application.Party("/follow").Handle(&controller.FollowController{})
			application.Party("/key").Handle(&controller.KeyController{})
		})

		//这里暂时先试用localAuth中间件作为验证本地 后续大概会添加签名或者其他验证机制 来确保可以多节点支持
		mvc.New(v0.Party("/", middleware.ApiLocalAuth)).Configure(func(application *mvc.Application) {
			application.Router.Post("/centerUser/createByPubKey", func(ctx iris.Context) {
				ctx.StopWithJSON(iris.StatusOK, new(controller.CenterUserController).CreateByPubKey(ctx))
			})
			application.Register(middleware.RegisterApiCenterUserAuth).Configure(func(application *mvc.Application) {
				application.Party("/centerUser").Handle(&controller.CenterUserController{})
			})
		})
	}
}
