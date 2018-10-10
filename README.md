# Go编写的自动部署

主要用于小项目托管在Gitee.com和Coding.net中的项目进行自动部署，将 `hook` 或者 `hook.exe` 运行，就可以运行一个接收hook的HTTP服务器

可多项目部署，如有多个项目，增加多个 `.json` 文件配置即可

#### 启动项目

    ./hook -port=8888

#### 访问路由

    Gitee：`/gitee`

    Coding: `/coding`

    例如：http://0.0.0.0:8888/gitee

### 文件说明

*  **git.sh**

    该文件是用于部署的，传的三个参数分别是 `项目在本机的目录` `项目的别名origin/master` `项目分支`

    例如：

    ```

     ./git.sh /home/wwwroot/abc origin/master master

    ```

*  **hook.go**

    源码


* **user.project.branch.json**

    ###### 例如你的coding用户名为test，项目名称为 himiweb，部署分支为 master， 则文件命名为： `test.himiweb.master.json`

    JSON配置文件，分别代指

    ```
    {
        "path": "项目在本机的目录，如 /home/wwwroot/abc",
        "head": "项目的别名，如origin/master",
        "password": "hook的验证密码"
    }
    ```

* **hook && hook.exe**

    两个文件分别是linux和windows下的可执行二进制文件
    ##### 如需更多系统版本，请自行运行 go build hook.go 打包
