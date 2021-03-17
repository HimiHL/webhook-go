# 前言

主要用于小项目托管在Gogs平台的项目进行自动部署，将 `hook`运行，就可以运行一个轻型的HTTP服务器

该项目主要是监听Gogs系统中，某个用户/组织下的所有项目，自动拉代码，无需再加配置文件

重要提示：如果修改了config.json 一定要reload


# 支持

- [x] 支持Linux平台下部署
- [x] 支持Git协议
- [x] 支持Http协议，需配置credential.helper
- [x] 支持自定义Shell脚本
- [x] 支持Gogs

# 命令列表

```shell
hook             # 启动服务器，-p 缺省，默认 `7442`
hook -p 7442     # 启动监听端口为7442的服务器

hook -h          # 帮助信息
hook -s reload   # 重启
hook -s stop     # 强行停止 SIGKILL
```

# 路由列表

* /

自动解析，将根据Header头中的参数自动解析，仅支持gogs

# 文件说明

*  **git.sh**

该文件是用于部署的，传的三个参数分别是 `项目在本机的目录` `项目分支` `项目的git地址（SSH协议/HTTP协议）`

也可以在**文件末尾**追加其他处理: 比如正式版发布时需要批量替换URL等操作

#### 可手动执行

```shell
./git.sh /home/wwwroot/abc master http://code.abc.com/xxx/xxx.git
```

*  **hook.go**

源码


* **config.json**
#### JSON配置解析

```json
[
    {
        "platform": "平台，目前仅支持 gogs",
        "namespace": "组织名称或者用户名称",
        "path": "代码存放的目录",
        "branch": "暂无用",
        "proto": "协议，http或者ssh",
        "password": "暂无用"
    }
]
```

* **hook**

Linux下可执行的二进制文件

* **log**
  
系统运行时生成的Log文件