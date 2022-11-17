package middleware

import (
	"fmt"
	"time"

	"github.com/ethtweet/ethtweet/config"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/models"

	"github.com/kataras/iris/v12"
)

func ApiLocalAuth(ctx iris.Context) {
	if !config.Cfg.CheckApiLocal || global.IsLocalIp(ctx.RemoteAddr()) || ctx.Method() == "GET" {
		ctx.Next()
		return
	}
	ctx.StopWithError(403, fmt.Errorf("invalid request method"))
	return
}

func RegisterApiCenterUserAuth(ctx iris.Context) *models.User {
	var ipfsHash, address string
	var err error
	if ctx.Method() == "GET" {
		ipfsHash = ctx.URLParamTrim("ipfsHash")
		address = ctx.URLParamTrim("address")
	} else {
		ipfsHash = ctx.PostValueTrim("ipfsHash")
		address = ctx.PostValueTrim("address")
	}
	user := &models.User{}
	if address != "" && global.GetDB().Where("id = ?", address).Limit(1).Find(user).RowsAffected > 0 {
		if user.IpfsHash == ipfsHash {
			ctx.Next()
			return user
		}
	}
	if ipfsHash == "" {
		ctx.StopWithError(403, fmt.Errorf("invalid ipfsHash"))
		return nil
	}
	//同步ipfs信息
	user2, err := models.GetUserByIpfs(ipfsHash)
	if err != nil || user2.Id != address {
		ctx.StopWithError(403, err)
		return nil
	}
	if user2.Nonce < user.Nonce {
		_ = user.UploadIpfs(global.GetDB())
		ctx.Next()
		return user
	}
	if user.Id == "" {
		user2.LocalNonce = 0
		user2.LocalUser = global.IsNo
		if err := global.GetDB().Create(user2).Error; err != nil {
			ctx.StopWithError(403, err)
			return nil
		}
		return user2
	}
	if user2.CreatedAt == 0 {
		user2.CreatedAt = time.Now().Unix()
	}
	//忽略掉指定字段
	if err = global.GetDB().Omit("LocalNonce", "LocalUser", "LatestCid").Save(user2).Error; err != nil {
		ctx.StopWithError(403, err)
		return nil
	}
	//将本地的数据覆盖ipfs上查询的返回
	user2.LocalUser = user.LocalUser
	user2.LocalNonce = user.LocalNonce
	user2.LatestCid = user.LatestCid
	//logs.PrintlnInfo("center api user ", user2)
	ctx.Next()
	return user2
}
