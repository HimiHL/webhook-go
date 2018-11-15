# 前言

主要用于小项目托管在[码云](http://gitee.com)和[Coding](http://coding.net)中的项目进行自动部署，将 `hook`  运行，就可以运行一个轻型的HTTP服务器

可多项目部署，如有多个项目，增加多个 `.json` 文件配置即可，__不需要重启hook进程__


# 支持

- [x] 支持Linux平台下部署
- [x] 支持Git协议，即: git@coding.net:user/project.git
- [ ] 不支持Http协议，即: https://git.coding.net/user/project.git
- [x] 支持自定义Shell脚本

# 命令列表

```shell
hook             # 启动服务器，-p 缺省，默认 `7442`
hook -p 7442     # 启动监听端口为7442的服务器

hook -h          # 帮助信息
hook -s reload   # 重启
hook -s stop     # 强行停止 SIGKILL
```

# 路由列表

* /gitee

解析[码云](http://gitee.com)的webhook通知，支持 `ContentType: application/json` 

* /coding
  
解析[Coding](http://coding.net)的webhook通知，支持V2版本的通知、`ContentType: application/json`


# 文件说明

*  **git.sh**

该文件是用于部署的，传的三个参数分别是 `项目在本机的目录` `项目的别名origin/master` `项目分支`

也可以在**文件末尾**追加其他处理: 比如正式版发布时需要批量替换URL等操作

#### 可手动执行

```shell
./git.sh /home/wwwroot/abc origin/master master
```

*  **hook.go**

源码


* **user.project.branch.json**

例如你的coding用户名为**test**，项目名称为 himiweb，部署分支为 master， 则文件命名为： `test.himiweb.master.json`

#### JSON配置解析

```json
{
    "path": "项目在本机的目录，如 /home/wwwroot/abc",
    "head": "项目的别名，如origin/master",
    "password": "hook的验证密码"
}
```

* **hook**

Linux下可执行的二进制文件

* **log**
  
系统运行时生成的Log文件