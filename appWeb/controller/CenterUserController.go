package controller

import (
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/ethtweet/ethtweet/appWeb"
	"github.com/ethtweet/ethtweet/broadcastMsg"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/keys"
	"github.com/ethtweet/ethtweet/logs"
	"github.com/ethtweet/ethtweet/models"
	"github.com/ethtweet/ethtweet/p2pNet"
	"strconv"
)

type CenterUserController struct {
	User *models.User
}

func (cu *CenterUserController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "updateInfo", "UpdateInfo")
	b.Handle("POST", "releaseTw", "ReleaseTw")
}

func (cu *CenterUserController) GetInfo(ctx iris.Context) *appWeb.ResponseFormat {
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", cu.User)
}

func (cu *CenterUserController) CreateByPubKey(ctx iris.Context) *appWeb.ResponseFormat {
	pubKey, err := global.Base58ToPubKey(ctx.PostValueTrim("pubKey"))
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}

	//通过公钥创建一个用户 该方法会为用户生成一个临时的peerId
	usr, err := models.GetOrCreateUserByPub(pubKey)
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	go func() {
		err = broadcastMsg.SyncUserInfo(usr, false)
		if err != nil {
			logs.PrintErr(err)
		}
	}()
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", map[string]interface{}{
		"address":  usr.Id,
		"name":     usr.Name,
		"ipfsHash": usr.IpfsHash,
	})
}

func (cu *CenterUserController) UpdateInfo(ctx iris.Context) *appWeb.ResponseFormat {
	nickname, errNickname := ctx.PostValueMany("nickname")
	avatar, errAvatar := ctx.PostValueMany("avatar")
	desc, errDesc := ctx.PostValueMany("desc")
	updateSignUnix, err := ctx.PostValueInt64("updateSignUnix")
	sign := ctx.PostValueTrim("sign")
	if err != nil || updateSignUnix <= cu.User.UpdatedSignUnix {
		return appWeb.NewResponse(appWeb.ResponseFailCode, fmt.Sprintf("invalid updateSignUnix %d <= %d", updateSignUnix, cu.User.UpdatedSignUnix), nil)
	}
	if sign == "" {
		return appWeb.NewResponse(appWeb.ResponseFailCode, "sign is empty", nil)
	}
	if (nickname == "" && avatar == "" && desc == "") || (nickname == cu.User.Name && avatar == cu.User.Avatar && desc == cu.User.Desc) && (sign == cu.User.Sign) {
		return appWeb.NewResponse(appWeb.ResponseFailCode, "invalid data", nil)
	}
	cu.User.UpdatedSignUnix = updateSignUnix

	if !errors.Is(errNickname, context.ErrNotFound) {
		cu.User.Name = nickname
	}
	if !errors.Is(errAvatar, context.ErrNotFound) {
		cu.User.Avatar = avatar
	}
	if !errors.Is(errDesc, context.ErrNotFound) {
		cu.User.Desc = desc
	}
	if !keys.VerifySignatureByAddress(cu.User.Id, sign, cu.User.GetSignMsg()) {
		return appWeb.NewResponse(appWeb.ResponseFailCode, fmt.Sprintf("sign is err %s, %s %s", sign, cu.User.GetSignMsg(), cu.User.Id), nil)
	}
	cu.User.Sign = sign
	if err := global.GetDB().Model(cu.User).Select("Name", "Avatar", "Desc", "UpdatedSignUnix", "Sign").Updates(cu.User).Error; err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}

	cUser := models.GetCurrentUser()
	if cUser != nil && cUser.UsrNode != nil {
		cu.User.UsrNode = models.GetCurrentUser().UsrNode
		//广播用户资料
		go func() {
			<-cu.User.UsrNode.WaitOnlineNode()
			cu.User.UsrNode.EachOnlineNodes(func(node *p2pNet.OnlineNode) bool {
				logs.PrintlnInfo("broadcast update info req to ", node.GetIdPretty())
				_ = p2pNet.WriteData(node.Rw, broadcastMsg.NewUserInfo(cu.User))
				return true
			})
		}()
	}
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", map[string]string{
		"ipfsHash": cu.User.IpfsHash,
	})
}

func (uc *CenterUserController) ReleaseTw(ctx iris.Context) *appWeb.ResponseFormat {
	tw := &models.Tweets{}
	tw.Id = ctx.PostValueTrim("id")
	tw.UserId = ctx.PostValueTrim("address")
	nonce, _ := ctx.PostValueInt64("nonce")
	tw.Nonce = uint64(nonce)
	tw.Content = ctx.PostValueTrim("content")
	tw.Attachment = ctx.PostValueTrim("attachment")
	tw.Sign = ctx.PostValueTrim("sign")
	tw.OriginTwId = ctx.PostValueTrim("origin_tw_id")
	tw.OriginUserId = ctx.PostValueTrim("origin_user_address")
	tw.TopicTag = ctx.PostValueTrim("topic_tag")
	t, err := strconv.ParseInt(ctx.PostValueTrim("created_at"), 10, 64)
	if err != nil {
		logs.PrintErr("need time ", err)
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	tw.CreatedAt = t
	logs.PrintlnInfo("center release tweets request start....", tw.Id)
	err = broadcastMsg.CenterUserRelease(tw)
	logs.PrintlnInfo("center release tweets request end....", tw.Id)
	if err != nil {
		logs.PrintErr("center release tweets request fail ", err)
		if errors.Is(err, global.ErrWaitUserSync) {
			go func() {
				_ = broadcastMsg.SyncUserInfo(uc.User, true)
			}()
		}
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", nil)
}
