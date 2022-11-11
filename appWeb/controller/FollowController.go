package controller

import (
	"github.com/ethtweet/ethtweet/appWeb"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/models"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type FollowController struct {
	User *models.User
}

func (fc *FollowController) BeforeActivation(b mvc.BeforeActivation) {}

func (fc *FollowController) Post(ctx iris.Context) *appWeb.ResponseFormat {
	fid := ctx.PostValueTrim("id")
	followed := &models.User{}
	global.GetDB().Model(followed).Where("id = ?", fid).Find(followed)

	if followed.Id == "" || followed.Id == fc.User.Id {
		return appWeb.NewResponse(appWeb.ResponseFailCode, "not found followed id", iris.Map{})
	}

	fl := models.Follow{
		UserId:     fc.User.Id,
		FollowedID: followed.Id,
	}

	global.GetDB().FirstOrCreate(&fl, fl)
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", fl)
}

func (fc *FollowController) PostCancel(ctx iris.Context) *appWeb.ResponseFormat {
	fid := ctx.PostValueTrim("id")
	followed := &models.User{}
	global.GetDB().Model(followed).Where("id = ?", fid).Find(followed)

	if followed.Id == "" {
		return appWeb.NewResponse(appWeb.ResponseFailCode, "not found followed id", iris.Map{})
	}

	fl := models.Follow{
		UserId:     fc.User.Id,
		FollowedID: followed.Id,
	}

	global.GetDB().Where(fl).Delete(&fl)
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", iris.Map{})
}

func (fc *FollowController) GetBy(id string, ctx iris.Context) *appWeb.ResponseFormat {
	user := &models.User{}
	global.GetDB().Model(user).Where("id = ?", id).Find(&user)

	if user.Id == "" {
		return appWeb.NewResponse(appWeb.ResponseFailCode, "not found", iris.Map{})
	}

	pager := global.NewPager(ctx)
	fls := make([]*models.Follow, 0, pager.Limit)
	global.GetDB().Where(models.Follow{UserId: user.Id}).Limit(pager.Limit).Order("created_at desc").Offset(pager.Offset).Find(&fls)
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", fls)
}

func (fc *FollowController) GetStateBy(id string, ctx iris.Context) *appWeb.ResponseFormat {
	fl := &models.Follow{}
	global.GetDB().Model(fl).Where(models.Follow{
		UserId:     fc.User.Id,
		FollowedID: id,
	}).Find(&fl)

	if fl.FollowedID != "" {
		return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", iris.Map{
			"State": 1,
		})
	}
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", iris.Map{
		"State": 0,
	})
}
