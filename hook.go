package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/bitly/go-simplejson"
	daemon "github.com/sevlyar/go-daemon"
)

var (
	help   = flag.Bool("h", false, "The Help")
	signal = flag.String("s", "", "send `signal` to a master process: stop, reload")
	port   = flag.String("p", "7442", "HTTP Server Port, Default `7442`")
)

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `Webhook version: webhook/1.0.0
	Usage: hook [-h] [-s signal] [-p port]

	Options: 
	`)
	flag.PrintDefaults()
}

func killHandler(sig os.Signal) error {
	log.Println("已停止运行")
	return nil
}

func reloadHandler(sig os.Signal) error {
	log.Println("服务器重载成功")
	return nil
}

/**
 * main主函数入口
 *
 */
func main() {

	// 命令行的参数获取
	flag.Parse()
	if *help {
		flag.Usage()
	} else {
		deamonHTTP()
	}

}

func deamonHTTP() {
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGKILL, killHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)
	cntxt := &daemon.Context{
		PidFileName: "pid",
		PidFilePerm: 0644,
		LogFileName: "log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        os.Args,
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalln("信号发送失败:", err)
		}
		daemon.SendCommands(d)
		return
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("无法启动服务器: ", err)
		return
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("进程运行")

	go serveHTTP()

	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}
	log.Println("daemon terminated")
}

func serveHTTP() {
	mux := http.NewServeMux()
	mux.HandleFunc("/gitee", gitee)
	mux.HandleFunc("/coding", coding)
	mux.HandleFunc("/", index)
	log.Fatalln(http.ListenAndServe(":"+*port, mux))
}

/**
 * gitee.com 的Webhook解析
 * 目前Content-Type只有JSON格式
 *
 */
func gitee(w http.ResponseWriter, request *http.Request) {
	contentType := request.Header.Get("Content-Type")
	if contentType == "application/json" {
		json := ParseGitEE(request)
		w.Write([]byte(json))
	} else {
		w.Write([]byte(`Hello GitEE`))
	}
}

/**
 * coding.net 的Webhook解析
 * 暂时不解析ContentType
 *
 */
func coding(w http.ResponseWriter, request *http.Request) {
	json := ParseCoding(request)
	w.Write([]byte(json))
}

/**
 *
 * 解析Coding.net的数据
 *
 */
func ParseCoding(request *http.Request) string {

	event := request.Header.Get(`X-Coding-Event`)

	if event == `ping` {
		return "这个Coding不简单"
	}

	result, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Print(`请求参数无法获取:` + err.Error())
		return "未获取到数据"
	}

	// 解析JSON
	json, err := simplejson.NewJson(result)
	if err != nil {
		log.Print(`JSON解析出错:` + err.Error())
		return "未获取到数据包"
	}

	hookID, err := json.Get(`hook_id`).String()
	if err == nil {
		return hookID + `一切正常`
	}

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
		log.Print(`无法读取文件:` + filename + `:` + err.Error())
		return "无法获取数据-"
	}

	fileJSON, err := simplejson.NewJson(b)
	if err != nil {
		log.Print(`JSON解析错误:` + err.Error())
		return "数据解析错误"
	}
	filePwd, _ := fileJSON.Get(`password`).String()
	filePath, _ := fileJSON.Get(`path`).String()
	fileHead, _ := fileJSON.Get(`head`).String()

	// 校验密码
	pwd, _ := json.Get(`token`).String()
	if pwd != `` {
		if pwd != filePwd {
			log.Print(`密码校验错误:` + pwd + `:正确密码:` + filePwd + `:` + err.Error())
			return "凭证校验异常"
		}
	}
	// 执行Shell 命令
	c := `./git.sh ` + filePath + ` ` + fileHead + ` ` + branch
	cmd := exec.Command("sh", "-c", c)
	err = cmd.Start() // 该操作不阻塞
	if err != nil {
		log.Print(`Shell执行异常:` + c + `:` + err.Error())
		return "任务执行异常"
	}
	return "Hello!"
}

/**
 * 解析Gitee.com
 *
 */
func ParseGitEE(request *http.Request) string {
	result, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Print(`请求参数无法获取:` + err.Error())
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
		log.Print(`无法读取文件:` + filename + `:` + err.Error())
		return "无法获取数据-"
	}

	fileJSON, err := simplejson.NewJson(b)
	if err != nil {
		log.Print(`JSON解析错误:` + err.Error())
		return "数据解析错误"
	}
	filePwd, _ := fileJSON.Get(`password`).String()
	filePath, _ := fileJSON.Get(`path`).String()
	fileHead, _ := fileJSON.Get(`head`).String()

	// 校验密码
	pwd, _ := json.Get(`password`).String()
	if pwd != `` {
		if pwd != filePwd {
			log.Print(`密码校验错误:` + pwd + `:正确密码:` + filePwd + `:` + err.Error())
			return "凭证校验异常"
		}
	}
	// 执行Shell 命令
	c := `./git.sh ` + filePath + ` ` + fileHead + ` ` + branch
	cmd := exec.Command("sh", "-c", c)
	// out, err := cmd.Output() // 该操作会阻塞
	err = cmd.Start() // 该操作不阻塞
	if err != nil {
		log.Print(`Shell执行异常:` + c + `:` + err.Error())
		return "任务执行异常"
	}
	return "Hello!"
}

func index(w http.ResponseWriter, request *http.Request) {
	w.Write([]byte(`Hello`))
}
