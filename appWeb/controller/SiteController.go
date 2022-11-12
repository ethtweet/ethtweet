package controller

import (
	"bytes"
	"github.com/ethtweet/ethtweet/appWeb"
	"github.com/ethtweet/ethtweet/global"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"io"
	"runtime"
)

type SiteController struct {
}

func (s *SiteController) BeforeActivation(b mvc.BeforeActivation) {}

func (s *SiteController) GetVersion(ctx iris.Context) *appWeb.ResponseFormat {
	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "OK", iris.Map{
		"golang":  runtime.Version(),
		"system":  runtime.GOARCH + "/" + runtime.GOOS,
		"version": global.Version,
	})
}

func (s *SiteController) PostUpload(ctx iris.Context) string {
	file, _, err := ctx.FormFile("up")
	if err != nil {
		return err.Error()
	}
	b, err := io.ReadAll(file)
	if err != nil {
		return err.Error()
	}
	hash, err := global.UploadIpfsReader(bytes.NewReader(b))
	if err != nil {
		return err.Error()
	}
	return hash
}
