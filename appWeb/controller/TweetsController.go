package controller

import (
	"errors"
	"fmt"
	"github.com/ethtweet/ethtweet/appWeb"
	"github.com/ethtweet/ethtweet/broadcastMsg"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/models"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type TweetsController struct {
	User *models.User
}

func (twc *TweetsController) BeforeActivation(b mvc.BeforeActivation) {}

func (twc *TweetsController) GetUserTimeline(ctx iris.Context) *appWeb.ResponseFormat {
	pager := global.NewPager(ctx)
	tws := make([]*models.Tweets, 0, pager.Limit)
	global.GetDB().Where("user_id", twc.User.Id).Limit(pager.Limit).
		Preload("OriginTw").
		Preload("UserInfo").
		Order("nonce desc").
		Offset(pager.Offset).Find(&tws)
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", tws)
}

func (twc *TweetsController) GetExplore(ctx iris.Context) *appWeb.ResponseFormat {
	pager := global.NewPager(ctx)
	tws := make([]*models.Tweets, 0, pager.Limit)
	global.GetDB().Limit(pager.Limit).Offset(pager.Offset).Preload("UserInfo").Order("created_at desc").Find(&tws)
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", tws)
}

func (twc *TweetsController) PostRelease(ctx iris.Context) *appWeb.ResponseFormat {
	keyName := ctx.PostValueTrim("key")
	var user *models.User
	if len(keyName) == 0 {
		keyName = twc.User.UsrNode.UserKey
	}

	user, err := models.GetUserByKeyName(twc.User.UsrNode.UserData, keyName, true)
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}

	createdAt, _ := ctx.PostValueInt64("createdAt")
	tw, err := broadcastMsg.ReleaseTweet(user,
		keyName, ctx.PostValueTrim("content"),
		ctx.PostValueTrim("attachment"),
		ctx.PostValueTrim("forward_id"),
		ctx.PostValueTrim("topic_tag"),
		createdAt)
	if err != nil {
		if errors.Is(err, global.ErrWaitUserSync) {
			go func() {
				_ = broadcastMsg.SyncUserInfo(twc.User, true)
			}()
		}
		return appWeb.NewResponse(appWeb.ResponseFailCode, fmt.Sprintf("release tweet err %s", err.Error()), err)
	}
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", tw)
}

func (twc *TweetsController) GetUserTimelineBy(id string, ctx iris.Context) *appWeb.ResponseFormat {
	pager := global.NewPager(ctx)
	tws := make([]*models.Tweets, 0, pager.Limit)
	global.GetDB().Where("user_id", id).Preload("OriginTw").Limit(pager.Limit).Order("nonce desc").Offset(pager.Offset).Preload("UserInfo").Find(&tws)
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", tws)
}

//我的转发
func (twc *TweetsController) GetUserForward(ctx iris.Context) *appWeb.ResponseFormat {
	pager := global.NewPager(ctx)
	tws := make([]*models.Tweets, 0, pager.Limit)
	global.GetDB().Where("user_id = ? and origin_user_id <> ''", twc.User.Id).Preload("OriginTw").Preload("UserInfo").Limit(pager.Limit).Offset(pager.Offset).Order("nonce desc").Find(&tws)
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", tws)
}

func (twc *TweetsController) GetAt_me(ctx iris.Context) *appWeb.ResponseFormat {
	pager := global.NewPager(ctx)
	tws := make([]*models.Tweets, 0, pager.Limit)
	global.GetDB().Where("origin_user_id = ?", twc.User.Id).Preload("OriginTw").Preload("UserInfo").Limit(pager.Limit).Offset(pager.Offset).Order("nonce desc").Find(&tws)
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", tws)
}

func (twc *TweetsController) GetBy(id string, ctx iris.Context) *appWeb.ResponseFormat {
	tw := models.Tweets{}
	global.GetDB().Where("id", id).Preload("UserInfo").Find(&tw)
	if tw.Id == "" {
		return appWeb.NewResponse(appWeb.ResponseFailCode, "not found", iris.Map{})
	}
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", tw)
}

func (twc *TweetsController) GetFollow(ctx iris.Context) *appWeb.ResponseFormat {
	pager := global.NewPager(ctx)
	tws := make([]*models.Tweets, 0, pager.Limit)

	fls := global.GetDB().Select("followed_id").Where("user_id", twc.User.Id).Table("follow")
	global.GetDB().
		Where("user_id IN (?)", fls).Limit(pager.Limit).
		Or("user_id", twc.User.Id).
		Order("created_at desc").Preload("UserInfo").Offset(pager.Offset).Find(&tws)

	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", tws)
}
