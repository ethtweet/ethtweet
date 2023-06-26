package controller

import (
	"github.com/ethtweet/ethtweet/appWeb"
	"github.com/ethtweet/ethtweet/broadcastMsg"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/logs"
	"github.com/ethtweet/ethtweet/models"
	"github.com/ethtweet/ethtweet/p2pNet"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type UserController struct {
	User *models.User
}

func (uc *UserController) BeforeActivation(b mvc.BeforeActivation) {}

func (uc *UserController) GetProfile(ctx iris.Context) *appWeb.ResponseFormat {
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", uc.User.GetUserInfoToPublic())
}

func (uc *UserController) GetNonceBy(id string, ctx iris.Context) *appWeb.ResponseFormat {
	user := &models.User{}
	if global.GetDB().Where("id", id).Limit(1).Find(user).RowsAffected == 0 {
		//询问用户
		go func() {
			uak := broadcastMsg.NewUserInfoAsk(&models.User{
				Id: id,
			})
			_ = uak.DoAsk()
		}()
		return appWeb.NewResponse(appWeb.ResponseFailCode, "not found user", nil)
	}
	var nonce uint64 = 0
	if user.Nonce > 0 || user.LatestCid != "" {
		nonce = user.Nonce + 1
	}

	var latestTw models.Tweets
	global.GetDB().Where(models.Tweets{UserId: id}).Order("nonce desc").First(&latestTw)
	if user.Nonce <= latestTw.Nonce {
		// fix user latest nonce
		global.GetDB().Model(&user).Update("nonce", latestTw.Nonce)
		nonce = latestTw.Nonce + 1
	}

	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", map[string]interface{}{
		"id":    id,
		"nonce": nonce,
	})
}

func (uc *UserController) GetBy(id string, ctx iris.Context) *appWeb.ResponseFormat {
	tws := models.User{}
	global.GetDB().Where("id", id).First(&tws)
	if tws.Id == "" {
		return appWeb.NewResponse(appWeb.ResponseFailCode, "not found", iris.Map{})
	}
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", tws.GetUserInfoToPublic())
}

func (uc *UserController) PostProfile(ctx iris.Context) *appWeb.ResponseFormat {

	name := ctx.PostValueTrim("name")
	desc := ctx.PostValueTrim("desc")
	avatar := ctx.PostValueTrim("avatar")

	if name == "" && desc == "" && avatar == "" {
		return appWeb.NewResponse(appWeb.ResponseFailCode, "nothing change", iris.Map{})
	}

	key := ctx.PostValueTrim("key")
	var user *models.User
	var err error
	if key == "" || key == uc.User.UsrNode.UserKey {
		user = uc.User
	} else {
		user, err = models.GetUserByKeyName(uc.User.UsrNode.UserData, key, true)
		if err != nil {
			return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), iris.Map{})
		}
	}

	if name != "" {
		user.Name = name
	}

	if desc != "" {
		user.Desc = desc
	}

	if avatar != "" {
		user.Avatar = avatar
	}
	user.UsrNode = uc.User.UsrNode
	user.Sign, err = user.UsrNode.SignMsg(key, user.GenerateSignMsg())
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), iris.Map{})
	}
	if err := global.GetDB().Save(user).Error; err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), iris.Map{})
	}
	go func() {
		<-user.UsrNode.WaitOnlineNode()
		user.UsrNode.EachOnlineNodes(func(node *p2pNet.OnlineNode) bool {
			logs.PrintlnInfo("broadcast update info req to ", node.Pi.ID)
			_ = p2pNet.WriteData(node.Rw, broadcastMsg.NewUserInfo(user))
			return true
		})
	}()

	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", user)
}
