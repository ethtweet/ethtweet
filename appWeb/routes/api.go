package routes

import (
	"github.com/ethtweet/ethtweet/appWeb/controller"
	"github.com/ethtweet/ethtweet/appWeb/middleware"
	"github.com/ethtweet/ethtweet/models"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func RegisterApiRoutes(app *iris.Application) {
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	app.Get("/webui", func(ctx iris.Context) {
		ctx.Redirect("https://ipfs.io/ipns/share.chaintweet.io")
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
