package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"webhook/model"

	daemon "github.com/sevlyar/go-daemon"
)

var (
	help   = flag.Bool("h", false, "The Help")
	signal = flag.String("s", "", "send `signal` to a master process: stop, reload")
	port   = flag.String("p", "7442", "HTTP Server Port, Default `7442`")
)

var appConfig model.Config

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
	loadConfig()
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

	loadConfig()

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("加载配置文件")

	go serveHTTP()

	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}
	log.Println("daemon terminated")
}

func serveHTTP() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", gogs)
	log.Fatalln(http.ListenAndServe(":"+*port, mux))
}

func gogs(w http.ResponseWriter, request *http.Request) {
	json := ParseGogs(request)
	w.Write([]byte(json))
}

func autoHook(configIndex int, owner string, name string, branch string, url string) string {

	// 读取文件【该项目的配置文件】
	filePath := fmt.Sprintf("%s/%s/%s/", appConfig[configIndex].Platform, owner, name)
	fileName := fmt.Sprintf("%s.json", branch)

	projectFileName := filePath + fileName

	var projectConfigModel model.ProjectConfig
	// 检查文件是否存在
	if !isExist(projectFileName) {
		log.Print(fmt.Sprintf("自动创建配置文件: %s", projectFileName))
		dirErr := CreateMutiDir(filePath)
		if dirErr != nil {
			log.Print(fmt.Sprintf("创建目录时出现错误: %s", dirErr.Error()))
			return "配置文件创建失败"
		}
		// 仓库存放的地方
		repositoryPath := fmt.Sprintf("%s/%s/%s/%s/%s", appConfig[0].Path, appConfig[configIndex].Platform, owner, name, branch)

		projectConfigModel.Path = repositoryPath
		projectConfigModel.Head = branch
		projectConfigModel.Password = appConfig[configIndex].Password
		projectConfigText, _ := json.Marshal(projectConfigModel)

		// 创建配置文件
		projectConfigF, _ := os.Create(projectFileName)
		io.WriteString(projectConfigF, string(projectConfigText))
		defer projectConfigF.Close()
	} else {
		projectConfigByte, err := ioutil.ReadFile(projectFileName)
		if err != nil {
			log.Println("读取项目配置文件出错" + err.Error())
			return "项目配置文件读取出错"
		}

		json.Unmarshal(projectConfigByte, &projectConfigModel)
	}

	// 检查目录是否存在
	if !isExist(projectConfigModel.Path) {
		CreateMutiDir(projectConfigModel.Path)
	}

	// 执行Shell 命令
	c := fmt.Sprintf("bash git.sh %s %s %s", projectConfigModel.Path, projectConfigModel.Head, url)
	log.Print(c)
	cmd := exec.Command("sh", "-c", c)
	err := cmd.Start() // 该操作不阻塞
	if err != nil {
		log.Print(`Shell执行异常:` + c + `:` + err.Error())
		return "任务执行异常"
	}
	return "The Job Done!"
}

func loadConfig() error {
	configFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Println("没有找到config.json文件")
		return err
	}
	json.Unmarshal(configFile, &appConfig)
	return nil
}

func ParseGogs(request *http.Request) string {

	result, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Print(`请求参数无法获取:` + err.Error())
		return "未获取到数据"
	}
	configIndex := 0
	var requestModel model.Gogs
	json.Unmarshal([]byte(result), &requestModel)

	// 拥有者的名称
	owner := requestModel.Repository.Owner.Username

	// 分支名称
	branch := strings.Split(requestModel.Ref, "/")[2]

	for index, config := range appConfig {
		if strings.EqualFold(config.Namespace, owner) && strings.EqualFold(config.Platform, "gogs") {
			configIndex = index
		}
	}

	if branch != appConfig[configIndex].Branch {
		return "俺不接受这个分支的push"
	}

	// 项目名称
	name := requestModel.Repository.Name

	// 仓库地址
	url := requestModel.Repository.SSHURL
	if strings.EqualFold(appConfig[configIndex].Proto, "http") {
		url = requestModel.Repository.CloneURL
	}

	return autoHook(configIndex, owner, name, branch, url)
}

//调用os.MkdirAll递归创建文件夹
func CreateMutiDir(filePath string) error {
	if !isExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			fmt.Println("创建文件夹失败,error info:", err)
			return err
		}
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func isExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
