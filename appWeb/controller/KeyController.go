package controller

import (
	"fmt"

	"github.com/ethtweet/ethtweet/appWeb"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/keys"
	"github.com/ethtweet/ethtweet/models"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	keystore "github.com/ipfs/boxo/keystore"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type KeyController struct {
	User *models.User
}

func (k *KeyController) BeforeActivation(b mvc.BeforeActivation) {}

// 导出私钥
func (k *KeyController) PostExport(ctx iris.Context) *appWeb.ResponseFormat {
	ks, err := keystore.NewFSKeystore(k.User.UsrNode.UserData)
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	pri, err := ks.Get(ctx.PostValueTrim("keyName"))
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	pri0, err := global.LibP2pPriToEthPri(pri)
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", hexutil.Encode(crypto.FromECDSA(pri0)))
}

// 列出key
func (k *KeyController) PostList() *appWeb.ResponseFormat {
	ks, err := keystore.NewFSKeystore(k.User.UsrNode.UserData)
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	list, err := ks.List()
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", list)
}

// 创建key
func (k *KeyController) PostGen(ctx iris.Context) *appWeb.ResponseFormat {
	ks, err := keystore.NewFSKeystore(k.User.UsrNode.UserData)
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	keyName := ctx.PostValueTrim("keyName")
	if ok, _ := ks.Has(keyName); ok {
		return appWeb.NewResponse(appWeb.ResponseFailCode, "key "+keyName+" already exists", nil)
	}
	priKey, err := keys.NewPrivateKey()
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	err = ks.Put(keyName, priKey.LibP2pPrivate)
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	_, _ = models.GetOrCreateUserByPri(priKey)
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", map[string]string{
		"Name": keyName,
	})
}

// 导入私钥
func (k *KeyController) PostImport(ctx iris.Context) *appWeb.ResponseFormat {
	var keyBase, keyName string
	if keyBase = ctx.PostValueTrim("keyBase"); keyBase == "" {
		return appWeb.NewResponse(appWeb.ResponseFailCode, "invalid key-base", nil)
	}
	if keyName = ctx.PostValueTrim("keyName"); keyName == "" {
		return appWeb.NewResponse(appWeb.ResponseFailCode, "invalid keyName", nil)
	}
	ks, err := keystore.NewFSKeystore(k.User.UsrNode.UserData)
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	if ok, _ := ks.Has(keyName); ok {
		return appWeb.NewResponse(appWeb.ResponseFailCode, "key "+keyName+" already exists", nil)
	}
	pri0, err := crypto.HexToECDSA(keyBase)
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	priKey, err := keys.NewPrivateKeyByEthPri(pri0)
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	err = ks.Put(keyName, priKey.LibP2pPrivate)
	_, _ = models.GetOrCreateUserByPri(priKey)
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", map[string]string{
		"Name": keyName,
	})
}

func (k *KeyController) PostRename(ctx iris.Context) *appWeb.ResponseFormat {
	var keyName, newKeyName string
	if keyName = ctx.PostValueTrim("keyName"); keyName == "" {
		return appWeb.NewResponse(appWeb.ResponseFailCode, "invalid keyName", nil)
	}
	if newKeyName = ctx.PostValueTrim("newKeyName"); newKeyName == "" {
		return appWeb.NewResponse(appWeb.ResponseFailCode, "invalid newKeyName", nil)
	}
	ks, err := keystore.NewFSKeystore(k.User.UsrNode.UserData)
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	if ok, _ := ks.Has(keyName); !ok {
		return appWeb.NewResponse(appWeb.ResponseFailCode, fmt.Sprintf("key name %s not exists", keyName), nil)
	}
	if ok, _ := ks.Has(newKeyName); ok {
		return appWeb.NewResponse(appWeb.ResponseFailCode, fmt.Sprintf("new key name %s already exists", newKeyName), nil)
	}
	pri, err := ks.Get(keyName)
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	err = ks.Put(newKeyName, pri)
	if err != nil {
		return appWeb.NewResponse(appWeb.ResponseFailCode, err.Error(), nil)
	}
	_ = ks.Delete(keyName)
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "success", map[string]string{
		"Was": keyName,
		"Now": newKeyName,
	})
}
