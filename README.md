<div align="center">
<h1>ChatGPT Dingtalk</h1>

[![Auth](https://img.shields.io/badge/Auth-eryajf-ff69b4)](https://github.com/eryajf)
[![Go Version](https://img.shields.io/github/go-mod/go-version/eryajf/chatgpt-dingtalk)](https://github.com/eryajf/chatgpt-dingtalk)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/eryajf/chatgpt-dingtalk)](https://github.com/eryajf/chatgpt-dingtalk/pulls)
[![GitHub Pull Requests](https://img.shields.io/github/stars/eryajf/chatgpt-dingtalk)](https://github.com/eryajf/chatgpt-dingtalk/stargazers)
[![HitCount](https://views.whatilearened.today/views/github/eryajf/chatgpt-dingtalk.svg)](https://github.com/eryajf/chatgpt-dingtalk)
[![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/eryajf/chatgpt-dingtalk)](https://hub.docker.com/r/eryajf/chatgpt-dingtalk)
[![Docker Pulls](https://img.shields.io/docker/pulls/eryajf/chatgpt-dingtalk)](https://hub.docker.com/r/eryajf/chatgpt-dingtalk)
[![GitHub license](https://img.shields.io/github/license/eryajf/chatgpt-dingtalk)](https://github.com/eryajf/chatgpt-dingtalk/blob/main/LICENSE)

<p> 🌉 在钉钉群聊中添加ChatGPT机器人 🌉</p>

<img src="https://camo.githubusercontent.com/82291b0fe831bfc6781e07fc5090cbd0a8b912bb8b8d4fec0696c881834f81ac/68747470733a2f2f70726f626f742e6d656469612f394575424971676170492e676966" width="800"  height="3">
</div><br>


## 前言

最近ChatGPT异常火爆，本项目可以助你将GPT机器人集成到钉钉群聊中。


> 🥳 **欢迎关注我的其他开源项目：**
>
> - [Go-Ldap-Admin](https://github.com/eryajf/go-ldap-admin)：🌉 基于Go+Vue实现的openLDAP后台管理项目。
> - [learning-weekly](https://github.com/eryajf/learning-weekly)：📝 周刊内容以运维技术和Go语言周边为主，辅以GitHub上优秀项目或他人优秀经验。
> - [HowToStartOpenSource](https://github.com/eryajf/HowToStartOpenSource)：🌈 GitHub开源项目维护协同指南。
> - [read-list](https://github.com/eryajf/read-list)：📖 优质内容订阅，阅读方为根本
> - [awesome-github-profile-readme-chinese](https://github.com/eryajf/awesome-github-profile-readme-chinese)：🦩 优秀的中文区个人主页搜集


## 功能简介

* 支持在钉钉群聊中添加机器人，通过@机器人进行聊天交互。
* 提问支持单聊与串聊两种模式，通过@机器人发关键字切换。

## 使用前提

* 有Openai账号，并且创建好`api_key`，注册相关事项可以参考[此文章](https://juejin.cn/post/7173447848292253704) 。访问[这里](https://beta.openai.com/account/api-keys)，申请个人秘钥。
* 在钉钉开发者后台创建机器人，配置应用程序回调。

## 使用教程

### 第一步，先创建机器人

创建步骤参考文档：[企业内部开发机器人](https://open.dingtalk.com/document/robots/enterprise-created-chatbot)，或者根据如下步骤进行配置。

1. 创建机器人。
   ![image_20221209_163616](https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20221209_163616.png)

   > `📢 注意：`可能现在创建机器人的时候名字为`chatgpt`会被钉钉限制，请用其他名字命名。

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
$ docker run -itd --name chatgpt -p 8090:8090 -e APIKEY=换成你的key -e SESSION_TIMEOUT=600 --restart=always  dockerproxy.com/eryajf/chatgpt-dingtalk:latest
```

运行命令中映射的配置文件参考下边的配置文件说明。

`第二种：基于配置文件挂载运行`

```sh
# 复制配置文件，根据自己实际情况，调整配置里的内容
$ cp config.dev.json config.json  # 其中 config.dev.json 从项目的根目录获取

# 运行项目
$ docker run -itd --name chatgpt -p 8090:8090  -v `pwd`/config.json:/app/config.json --restart=always  dockerproxy.com/eryajf/chatgpt-dingtalk:latest
```

其中配置文件参考下边的配置文件说明。

注意，不论通过上边哪种docker方式部署，都需要配置Nginx代理，当然你直接通过服务器外网IP也可以。

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

Nginx配置完毕之后，可以先手动请求一下，通过服务日志输出判断服务是否正常可用：

```sh
$ curl --location --request POST 'http://chat.eryajf.net/' \
  --header 'Content-type: application/json' \
  --data-raw '{
    "conversationId": "xxx",
    "atUsers": [
        {
            "dingtalkId": "xxx",
            "staffId":"xxx"
        }
    ],
    "chatbotCorpId": "dinge8a565xxxx",
    "chatbotUserId": "$:LWCP_v1:$Cxxxxx",
    "msgId": "msg0xxxxx",
    "senderNick": "eryajf",
    "isAdmin": true,
    "senderStaffId": "user123",
    "sessionWebhookExpiredTime": 1613635652738,
    "createAt": 1613630252678,
    "senderCorpId": "dinge8a565xxxx",
    "conversationType": "2",
    "senderId": "$:LWCP_v1:$Ff09GIxxxxx",
    "conversationTitle": "机器人测试-TEST",
    "isInAtList": true,
    "sessionWebhook": "https://oapi.dingtalk.com/robot/sendBySession?session=xxxxx",
    "text": {
        "content": " 你好"
    },
    "msgtype": "text"
}'
```

如果手动请求没有问题，那么就可以在钉钉群里与机器人进行对话了。

`帮助列表`

> 艾特机器人发送空内容或者帮助，会返回帮助列表。

![image_20230215_184402](https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230215_184402.png)

`切换模式`

> 发送指定关键字，可以切换不同的模式。

![image_20230215_184655](https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230215_184655.png)

> 📢 注意：串聊模式下，群里每个人的聊天上下文是独立的。
> 📢 注意：默认对话模式为单聊，因此不必发送单聊即可进入单聊模式，而要进入串聊，则需要发送串聊关键字进行切换，当串聊内容超过最大限制的时候，你可以发送重置，然后再次进入串聊模式。

`实际聊天效果如下`

![image_20221209_163739](https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20221209_163739.png)

---

如果你想通过命令行直接部署，可以直接下载release中的[压缩包](https://github.com/eryajf/chatgpt-dingtalk/releases) ，请根据自己系统以及架构选择合适的压缩包，下载之后直接解压运行。

下载之后，在本地解压，即可看到可执行程序，与配置文件：

```
$ tar xf chatgpt-dingtalk-v0.0.4-darwin-arm64.tar.gz
$ cd chatgpt-dingtalk-v0.0.4-darwin-arm64
$ cp config.dev.json # 根据情况调整配置文件内容
$ ./chatgpt-dingtalk  # 直接运行

# 如果要守护在后台运行
$ nohup ./chatgpt-dingtalk &> run.log &
$ tail -f run.log
```


## 本地开发

```sh
# 获取项目
$ git clone https://github.com/eryajf/chatgpt-dingtalk.git

# 进入项目目录
$ cd chatgpt-dingtalk

# 复制配置文件，根据个人实际情况进行配置
$ cp config.dev.json config.json

# 启动项目
$ go run main.go
```

## 配置文件说明

```json
{
    "api_key": "xxxxxxxxx",  // openai api_key
    "session_timeout": 600   // 会话超时时间,默认600秒,在会话时间内所有发送给机器人的信息会作为上下文
}
```

## 常见问题

一些常见的问题，我单独开issue放在这里：[点我](https://github.com/eryajf/chatgpt-dingtalk/issues/44)，可以查看这里辅助你解决问题，如果里边没有，请对历史issue进行搜索(不要提交重复的issue)，也欢迎大家补充。

## 高光时刻

> 本项目曾在[2022-12-12](https://github.com/bonfy/github-trending/blob/master/2022/2022-12-12.md#go),[2022-12-18](https://github.com/bonfy/github-trending/blob/master/2022/2022-12-18.md#go),[2022-12-19](https://github.com/bonfy/github-trending/blob/master/2022/2022-12-19.md#go),[2022-12-20](https://github.com/bonfy/github-trending/blob/master/2022/2022-12-20.md#go),[2023-02-09](https://github.com/bonfy/github-trending/blob/master/2023-02-09.md#go),[2023-02-10](https://github.com/bonfy/github-trending/blob/master/2023-02-10.md#go),[2023-02-11](https://github.com/bonfy/github-trending/blob/master/2023-02-11.md#go),[2023-02-12](https://github.com/bonfy/github-trending/blob/master/2023-02-12.md#go)，这些天里，登上GitHub Trending。而且还在持续登榜中，可见最近openai的热度。
> ![image_20230215_094034](https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230215_094034.png)