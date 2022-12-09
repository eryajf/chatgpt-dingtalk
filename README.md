<div align="center">
<h1>Chatgpt Dingtalk</h1>

[![Auth](https://img.shields.io/badge/Auth-eryajf-ff69b4)](https://github.com/eryajf)
[![Go Version](https://img.shields.io/github/go-mod/go-version/eryajf/chatgpt-dingtalk)](https://github.com/eryajf/chatgpt-dingtalk)
[![Gin Version](https://img.shields.io/badge/Gin-1.6.3-brightgreen)](https://github.com/eryajf/chatgpt-dingtalk)
[![Gorm Version](https://img.shields.io/badge/Gorm-1.20.12-brightgreen)](https://github.com/eryajf/chatgpt-dingtalk)
[![GitHub Issues](https://img.shields.io/github/issues/eryajf/chatgpt-dingtalk.svg)](https://github.com/eryajf/chatgpt-dingtalk/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/eryajf/chatgpt-dingtalk)](https://github.com/eryajf/chatgpt-dingtalk/pulls)
[![GitHub Pull Requests](https://img.shields.io/github/stars/eryajf/chatgpt-dingtalk)](https://github.com/eryajf/chatgpt-dingtalk/stargazers)
[![HitCount](https://views.whatilearened.today/views/github/eryajf/chatgpt-dingtalk.svg)](https://github.com/eryajf/chatgpt-dingtalk)
[![GitHub license](https://img.shields.io/github/license/eryajf/chatgpt-dingtalk)](https://github.com/eryajf/chatgpt-dingtalk/blob/main/LICENSE)

<p> 🌉 在钉钉群聊中添加chatGPT机器人 🌉</p>

<img src="https://camo.githubusercontent.com/82291b0fe831bfc6781e07fc5090cbd0a8b912bb8b8d4fec0696c881834f81ac/68747470733a2f2f70726f626f742e6d656469612f394575424971676170492e676966" width="800"  height="3">
</div><br>


## 前言

最近chatGPT异常火爆，本项目可以将GPT机器人集成到钉钉群聊中。

`感谢：`这个项目借鉴了[wechatbot](https://github.com/869413421/wechatbot.git)，wechatbot是一个能够集成到个人微信的GPT机器人。

### 功能简介

* 支持在钉钉群聊中添加机器人，通过@机器人进行聊天交互。
* 提问增加上下文(可能不太理想)，更接近官网效果。

## 使用前提

* 有openai账号，并且创建好api_key，注册相关事项可以参考[此文章](https://juejin.cn/post/7173447848292253704) 。访问[这里](https://beta.openai.com/account/api-keys)，申请个人秘钥。
* 在钉钉开发者后台创建机器人，配置应用程序回调。

## 使用教程

### 第一步，先创建机器人

创建步骤参考文档：[企业内部开发机器人](https://open.dingtalk.com/document/robots/enterprise-created-chatbot)，或者根据如下步骤进行配置。

1. 创建机器人。
   ![image_20221209_163616](https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20221209_163616.png)

   步骤比较简单，这里就不赘述了。

2. 配置机器人回调接口。
   ![image_20221209_163652](https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20221209_163652.png)

   创建完毕之后，点击机器人开发管理，然后配置将要部署的服务所在服务器的出口IP，以及将要给服务配置的域名。

3. 发布机器人。
   ![image_20221209_163709](https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20221209_163709.png)

   点击版本管理与发布，然后点击上线，这个时候就能在钉钉的群里中添加这个机器人了。

4. 群聊添加机器人。

   ![image_20221209_163724](https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20221209_163724.png)

### 第二步，部署应用

你可以使用docker快速运行本项目。

`第一种：基于环境变量运行`

```sh
# 运行项目
$ docker run -itd --name chatgpt -e ApiKey=xxxx -e SessionTimeout=60s --restart=always docker.mirrors.sjtug.sjtu.edu.cn/eryajf/chatgpt-dingtalk:latest
```

运行命令中映射的配置文件参考下边的配置文件说明。

`第二种：基于配置文件挂载运行`

```sh
# 复制配置文件，根据自己实际情况，调整配置里的内容
$ cp config.dev.json config.json  # 其中 config.dev.json 从项目的根目录获取

# 运行项目
docker run -itd --name chatgpt -v ./config.json:/app/config.json --restart=always docker.mirrors.sjtug.sjtu.edu.cn/eryajf/chatgpt-dingtalk:latest
```

其中配置文件参考下边的配置文件说明。

部署完成之后，通过Nginx代理本服务：

```nginx
server {
    listen       80;
    server_name  chat.eryajf.net;

    client_header_timeout 120s;
    client_body_timeout 120s;

    location / {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_pass http://localhost:8090;
    }
}
```

部署完成之后，就可以在群里艾特机器人进行体验了。

效果如下：

![image_20221209_163739](https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20221209_163739.png)

## 本地开发

````sh
# 获取项目
$ git clone https://github.com/eryajf/chatgpt-dingtalk.git

# 进入项目目录
$ cd chatgpt-dingtalk

# 复制配置文件，根据个人实际情况进行配置
$ cp config.dev.json config.json

# 启动项目
$ go run main.go
````

## 配置文件说明
````json
{
    "api_key": "xxxxxxxxx",  // openai api_key
    "session_timeout": 60    // 会话超时时间,默认60秒,在会话时间内所有发送给机器人的信息会作为上下文
}
````