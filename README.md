# Go编写的自动部署

> 启动项目

```
./linux_hook -d=true
```

> 开启防火墙端口 

`7442`


> 访问地址

Gitee.com的项目为：`ip:7442/gitee`

待增加

### 文件说明


> `git.sh`

该文件是用于部署的，传的三个参数分别是 `项目在本机的目录` `项目的别名origin/master` `项目分支`

> `hook.go` 

源码

> `youngrain.xixiweb.master.json`

JSON配置文件，分别代指

```
{
    "path": "项目在本机的目录",
    "head": "项目的别名origin/master",
    "password": "hook的验证密码"
}
```


##### 如需更多系统版本，请自行运行 go build hook.go 打包

代码写的烂，不接受批评