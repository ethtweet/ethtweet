package update

import (
	"archive/zip"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/logs"
	"github.com/polydawn/refmt/json"
)

func ChcckGithubVersion() {
	r, err := http.Get("https://api.github.com/repos/ethtweet/ethtweet/releases/latest")
	if err != nil {
		return
	}
	b, err := io.ReadAll(r.Body)
	var v interface{}
	err = json.Unmarshal(b, &v)
	if err != nil {
		logs.PrintErr(err)
		return
	}

	data := v.(map[string]interface{})

	githubVerion := fmt.Sprintf("%s", data["tag_name"])
	githubVerion = strings.Replace(githubVerion, "v", "", 1)
	if compareVersion(githubVerion, global.Version) > 0 {
		logs.PrintlnSuccess("GitHub版本更高")
	} else {
		logs.PrintlnSuccess("不需要升级")
		return
	}

	githubPublishedTime, _ := time.ParseInLocation("2006-01-02T15:04:05Z", fmt.Sprintf("%s", data["published_at"]), time.Local)
	if time.Now().Sub(githubPublishedTime) < (time.Second * 3600) {
		logs.PrintlnSuccess("更新时间不足1个小时，延迟更新")
		return
	}
	updateFileUrl := fmt.Sprintf("https://github.com/ethtweet/ethtweet/releases/download/v%s/EthTweet-%s-%s-%s.zip", githubVerion, githubVerion, runtime.GOOS, runtime.GOARCH)
	// Get the data
	resp, err := http.Get(updateFileUrl)

	if resp.StatusCode != 404 {
		logs.PrintErr("文件不存在，404错误")
		return
	}
	if err != nil {
		logs.PrintErr(err)
		return
	} else {
		logs.PrintlnSuccess("下载最新安装包成功")
	}
	defer resp.Body.Close()

	// 创建一个文件用于保存
	out, err := os.Create("update.zip")
	if err != nil {
		logs.PrintErr(err)
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logs.PrintErr(err)
		return
	}

	h := sha512.New()
	if _, err := io.Copy(h, out); err != nil {
		logs.PrintErr(err)
		return
	}

	fileSha512 := hex.EncodeToString(h.Sum(nil))

	checksumsFileURL := fmt.Sprintf("https://github.com/ethtweet/ethtweet/releases/download/v%s/EthTweet-%s-%s-%s.zip.sha512", githubVerion, githubVerion, runtime.GOOS, runtime.GOARCH)
	r, err = http.Get(checksumsFileURL)
	if err != nil {
		logs.PrintErr(err)
		return
	}
	b, err = io.ReadAll(r.Body)
	checksums := string(b)
	if strings.Index(checksums, fileSha512) < 0 {

		logs.PrintErr("文件sha512错误")
		return
	}

	exeFilename, _ := os.Executable()

	//删除老文件
	if global.FileExists(path.Base(exeFilename) + ".old") {
		err = os.Remove(path.Base(exeFilename) + ".old")
		if err != nil {
			logs.PrintErr(err)
			return
		}
	}

	err = os.Rename(path.Base(exeFilename), path.Base(exeFilename)+".old")
	if err != nil {
		logs.PrintErr(err)
		return
	}

	err = Unzip("update.zip", ".")
	if err != nil {
		logs.PrintErr(err)
		return
	}

	logs.Println("current version: ", global.Version)
	logs.Println("Update to version: ", githubVerion)
	logs.Println("Ready to restart")
	time.Sleep(time.Second * 5) //更新前休眠5秒，避免重复冲突
	os.Exit(0)
}

func Unzip(zipPath, dstDir string) error {
	// open zip file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		if err := unzipFile(file, dstDir); err != nil {
			return err
		}
	}
	return nil
}

func unzipFile(file *zip.File, dstDir string) error {
	// create the directory of file
	filePath := path.Join(dstDir, file.Name)
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// open the file
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// create the file
	w, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer w.Close()

	w.Chmod(0777)

	// save the decompressed file content
	_, err = io.Copy(w, rc)
	return err
}

func compareVersion(version1 string, version2 string) int {
	var res int
	ver1Strs := strings.Split(version1, ".")
	ver2Strs := strings.Split(version2, ".")
	ver1Len := len(ver1Strs)
	ver2Len := len(ver2Strs)
	verLen := ver1Len
	if len(ver1Strs) < len(ver2Strs) {
		verLen = ver2Len
	}
	for i := 0; i < verLen; i++ {
		var ver1Int, ver2Int int
		if i < ver1Len {
			ver1Int, _ = strconv.Atoi(ver1Strs[i])
		}
		if i < ver2Len {
			ver2Int, _ = strconv.Atoi(ver2Strs[i])
		}
		if ver1Int < ver2Int {
			res = -1
			break
		}
		if ver1Int > ver2Int {
			res = 1
			break
		}
	}
	return res
}
