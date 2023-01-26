package controller

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/ethtweet/ethtweet/appWeb"
	"github.com/ethtweet/ethtweet/global"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
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

func (s *SiteController) GetBootstrap(ctx iris.Context) *appWeb.ResponseFormat {

	var bootstraps []string
	if global.FileExists("Bootstrap.txt") {
		file2, err := os.Open("Bootstrap.txt")
		if err != nil {
			return appWeb.NewResponse(appWeb.ResponseSuccessCode, "OK", iris.Map{
				"bootstraps": bootstraps,
			})
		}
		reader := bufio.NewReader(file2)
		for {
			str, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			bootstraps = append(bootstraps, strings.Trim(str, "\n"))

		}
	}

	return appWeb.NewResponse(appWeb.ResponseSuccessCode, "OK", iris.Map{
		"bootstraps": bootstraps,
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
