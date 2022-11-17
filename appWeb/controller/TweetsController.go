package controller

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ethtweet/ethtweet/appWeb"
	"github.com/ethtweet/ethtweet/broadcastMsg"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/keys"
	"github.com/ethtweet/ethtweet/logs"
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

// 广播签名后的推文
func (twc *TweetsController) PostReleaseBySign(ctx iris.Context) *appWeb.ResponseFormat {
	tw := &models.Tweets{}
	tw.UserId = ctx.PostValueTrim("address")
	nonce, _ := ctx.PostValueInt64("nonce")
	tw.Nonce = uint64(nonce)
	tw.Content = ctx.PostValueTrim("content")
	tw.Attachment = ctx.PostValueTrim("attachment")
	tw.Sign = ctx.PostValueTrim("sign")
	tw.OriginTwId = ctx.PostValueTrim("origin_tw_id")
	tw.OriginUserId = ctx.PostValueTrim("origin_user_address")
	tw.TopicTag = ctx.PostValueTrim("topic_tag")
	tw.GenerateId()
	t, err := strconv.ParseInt(ctx.PostValueTrim("created_at"), 10, 64)
	if err != nil {
		logs.PrintErr("need time ", err)
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	tw.CreatedAt = t
	if !keys.VerifySignatureByAddress(tw.UserId, tw.Sign, tw.GetSignMsg()) {
		logs.PrintErr("tw sign err %s %s %s", tw.UserId, tw.Sign, tw.GetSignMsg())
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	err = broadcastMsg.CenterUserRelease(tw)
	go broadcastMsg.BroadcastTweetSync(tw)

	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", tw)
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
				_ = broadcastMsg.SyncUserInfo(user, true)
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

// 我的转发
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
