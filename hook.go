package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	} else {
		w.Write([]byte(`Hello GitEE`))
	}
}

// 解析Coding的请求
func coding(w http.ResponseWriter, request *http.Request) {
	// 获取响应头
	contentType := request.Header.Get("Content-Type")
	if contentType == "application/json" {
		json := ParseCoding(request)
		w.Write([]byte(json))
	} else {
		w.Write([]byte(`Hello Coding`))
	}
}

/**
解析JSON数据
 *
*/
func ParseCoding(request *http.Request) string {
	result, err := ioutil.ReadAll(request.Body)
	if err != nil {
		Logger(`请求参数无法获取:` + err.Error())
		return "未获取到数据"
	}

	// 解析JSON
	json, err := simplejson.NewJson(result)

	// 分支名称
	ref, _ := json.Get(`ref`).String()
	branchs := strings.Split(ref, `/`)
	branch := branchs[2]

	// 获取拥有者
	owner, _ := json.Get(`repository`).Get(`owner`).Get(`name`).String()
	projectName, _ := json.Get(`repository`).Get(`name`).String()

	// 读取文件
	filename := owner + `.` + projectName + `.` + branch + `.json`

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		Logger(`无法读取文件:` + filename + `:` + err.Error())
		return "无法获取数据-"
	}

	fileJSON, err := simplejson.NewJson(b)
	if err != nil {
		Logger(`JSON解析错误:` + err.Error())
		return "数据解析错误"
	}
	filePwd, _ := fileJSON.Get(`password`).String()
	filePath, _ := fileJSON.Get(`path`).String()
	fileHead, _ := fileJSON.Get(`head`).String()

	// 校验密码
	pwd, _ := json.Get(`token`).String()
	if pwd != `` {
		if pwd != filePwd {
			Logger(`密码校验错误:` + pwd + `:正确密码:` + filePwd + `:` + err.Error())
			return "凭证校验异常"
		}
	}
	// 执行Shell 命令
	c := `./git.sh ` + filePath + ` ` + fileHead + ` ` + branch
	cmd := exec.Command("sh", "-c", c)
	// out, err := cmd.Output() // 该操作会阻塞
	err = cmd.Start() // 该操作不阻塞
	if err != nil {
		Logger(`Shell执行异常:` + c + `:` + err.Error())
		return "任务执行异常"
	}
	return "Hello!"
}

/**
解析JSON数据
 *
*/
func ParseGitEE(request *http.Request) string {
	result, err := ioutil.ReadAll(request.Body)
	if err != nil {
		Logger(`请求参数无法获取:` + err.Error())
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
	filename := projectName[0] + `.` + projectName[1] + `.` + branch + `.json`

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		Logger(`无法读取文件:` + filename + `:` + err.Error())
		return "无法获取数据-"
	}
	return string(b)

	fileJSON, err := simplejson.NewJson(b)
	if err != nil {
		Logger(`JSON解析错误:` + err.Error())
		return "数据解析错误"
	}
	filePwd, _ := fileJSON.Get(`password`).String()
	filePath, _ := fileJSON.Get(`path`).String()
	fileHead, _ := fileJSON.Get(`head`).String()

	// 校验密码
	pwd, _ := json.Get(`password`).String()
	if pwd != `` {
		if pwd != filePwd {
			Logger(`密码校验错误:` + pwd + `:正确密码:` + filePwd + `:` + err.Error())
			return "凭证校验异常"
		}
	}
	// 执行Shell 命令
	c := `./git.sh ` + filePath + ` ` + fileHead + ` ` + branch
	cmd := exec.Command("sh", "-c", c)
	// out, err := cmd.Output() // 该操作会阻塞
	err = cmd.Start() // 该操作不阻塞
	if err != nil {
		Logger(`Shell执行异常:` + c + `:` + err.Error())
		return "任务执行异常"
	}
	return "Hello!"
}

func index(w http.ResponseWriter, request *http.Request) {
	w.Write([]byte(`Hello`))
}

func Logger(message string) {
	filename := "runtime.log"
	logFile, err := os.Create(filename)
	defer logFile.Close()
	if err != nil {
	}
	debugLog := log.New(logFile, "[Info]", log.Llongfile)
	debugLog.Println(message)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/gitee", gitee)
	mux.HandleFunc("/coding", coding)
	mux.HandleFunc("/", index)
	log.Fatalln(http.ListenAndServe(":7442", mux))
}
