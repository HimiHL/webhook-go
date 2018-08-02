package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/bitly/go-simplejson"
	_ "github.com/icattlecoder/godaemon"
)

// 解析GitEE的请求
func gitee(w http.ResponseWriter, request *http.Request) {
	// 获取响应头
	contentType := request.Header.Get("Content-Type")
	if contentType == "application/json" {
		json := ParseGitEE(request)
		w.Write([]byte(json))
	}

}

/**
解析JSON数据
 *
*/
func ParseGitEE(request *http.Request) string {
	result, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return "未获取到数据"
	}

	// 解析JSON
	json, err := simplejson.NewJson(result)

	// 分支名称
	ref, _ := json.Get(`ref`).String()
	branchs := strings.Split(ref, `/`)
	branch := branchs[2]

	// 获取项目名称
	projName, _ := json.Get(`repository`).Get(`path_with_namespace`).String()
	projectName := strings.Split(projName, `/`)

	// 读取文件
	filename := projectName[0] + `_` + projectName[1] + `_` + branch + `.json`

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "无法获取数据-"
	}

	fileJSON, err := simplejson.NewJson(b)
	if err != nil {
		return "数据解析错误"
	}
	filePwd, _ := fileJSON.Get(`password`).String()
	filePath, _ := fileJSON.Get(`path`).String()
	fileHead, _ := fileJSON.Get(`head`).String()

	// 校验密码
	pwd, _ := json.Get(`password`).String()
	if pwd != filePwd {
		return "凭证校验异常"
	}
	// 执行Shell 命令
	c := `git.sh ` + filePath + ` ` + fileHead + ` ` + branch
	cmd := exec.Command("sh", "-c", c)
	out, err := cmd.Output()
	if err != nil {
		return "任务执行异常"
	}
	fmt.Println(out)
	return "Hello!"
}

func index(w http.ResponseWriter, request *http.Request) {
	w.Write([]byte(`Hello`))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/gitee", gitee)
	mux.HandleFunc("/", index)
	log.Fatalln(http.ListenAndServe(":7442", mux))
}
