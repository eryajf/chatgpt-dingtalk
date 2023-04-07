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

<p> 🌉 基于GO语言实现的钉钉集成ChatGPT机器人 🌉</p>

<img src="https://camo.githubusercontent.com/82291b0fe831bfc6781e07fc5090cbd0a8b912bb8b8d4fec0696c881834f81ac/68747470733a2f2f70726f626f742e6d656469612f394575424971676170492e676966" width="800"  height="3">
</div><br>

![](https://user-images.githubusercontent.com/33259379/223607306-2ac836a2-7ce5-4a12-a16e-bec40b22d8d6.png)


<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**目录**

- [前言](#%E5%89%8D%E8%A8%80)
- [功能介绍](#%E5%8A%9F%E8%83%BD%E4%BB%8B%E7%BB%8D)
- [使用前提](#%E4%BD%BF%E7%94%A8%E5%89%8D%E6%8F%90)
- [使用教程](#%E4%BD%BF%E7%94%A8%E6%95%99%E7%A8%8B)
  - [第一步，创建机器人](#%E7%AC%AC%E4%B8%80%E6%AD%A5%E5%88%9B%E5%BB%BA%E6%9C%BA%E5%99%A8%E4%BA%BA)
    - [方案一：outgoing类型机器人](#%E6%96%B9%E6%A1%88%E4%B8%80outgoing%E7%B1%BB%E5%9E%8B%E6%9C%BA%E5%99%A8%E4%BA%BA)
    - [方案二：企业内部应用](#%E6%96%B9%E6%A1%88%E4%BA%8C%E4%BC%81%E4%B8%9A%E5%86%85%E9%83%A8%E5%BA%94%E7%94%A8)
  - [第二步，部署应用](#%E7%AC%AC%E4%BA%8C%E6%AD%A5%E9%83%A8%E7%BD%B2%E5%BA%94%E7%94%A8)
    - [docker部署](#docker%E9%83%A8%E7%BD%B2)
    - [二进制部署](#%E4%BA%8C%E8%BF%9B%E5%88%B6%E9%83%A8%E7%BD%B2)
- [亮点特色](#%E4%BA%AE%E7%82%B9%E7%89%B9%E8%89%B2)
  - [与机器人私聊](#%E4%B8%8E%E6%9C%BA%E5%99%A8%E4%BA%BA%E7%A7%81%E8%81%8A)
  - [帮助列表](#%E5%B8%AE%E5%8A%A9%E5%88%97%E8%A1%A8)
  - [切换模式](#%E5%88%87%E6%8D%A2%E6%A8%A1%E5%BC%8F)
  - [查询余额](#%E6%9F%A5%E8%AF%A2%E4%BD%99%E9%A2%9D)
  - [日常问题](#%E6%97%A5%E5%B8%B8%E9%97%AE%E9%A2%98)
  - [通过内置prompt聊天](#%E9%80%9A%E8%BF%87%E5%86%85%E7%BD%AEprompt%E8%81%8A%E5%A4%A9)
  - [生成图片](#%E7%94%9F%E6%88%90%E5%9B%BE%E7%89%87)
  - [支持 gpt-4](#%E6%94%AF%E6%8C%81-gpt-4)
- [本地开发](#%E6%9C%AC%E5%9C%B0%E5%BC%80%E5%8F%91)
- [配置文件说明](#%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E8%AF%B4%E6%98%8E)
- [常见问题](#%E5%B8%B8%E8%A7%81%E9%97%AE%E9%A2%98)
- [进群交流](#%E8%BF%9B%E7%BE%A4%E4%BA%A4%E6%B5%81)
- [感谢](#%E6%84%9F%E8%B0%A2)
- [赞赏](#%E8%B5%9E%E8%B5%8F)
- [高光时刻](#%E9%AB%98%E5%85%89%E6%97%B6%E5%88%BB)
- [贡献者列表](#%E8%B4%A1%E7%8C%AE%E8%80%85%E5%88%97%E8%A1%A8)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 前言

本项目可以助你将GPT机器人集成到钉钉群聊当中。当前默认模型为`gpt-3.5`，支持`gpt-4`。

>- `📢 注意`：当下部署以及配置流程都已非常成熟，文档和issue中基本都覆盖到了，因此不再回答任何项目安装部署与配置使用上的问题，如果完全不懂，可考虑通过 **[邮箱](mailto:eryajf@163.com)** 联系我进行付费的技术支持。
>
>- `📢 注意`：这个项目所有的功能，都汇聚在[使用指南](./docs/userGuide.md)中，请务必仔细阅读，以体验其完整精髓。

🥳 **欢迎关注我的其他开源项目：**
>
> - [Go-Ldap-Admin](https://github.com/eryajf/go-ldap-admin)：🌉 基于Go+Vue实现的openLDAP后台管理项目。
> - [learning-weekly](https://github.com/eryajf/learning-weekly)：📝 周刊内容以运维技术和Go语言周边为主，辅以GitHub上优秀项目或他人优秀经验。
> - [HowToStartOpenSource](https://github.com/eryajf/HowToStartOpenSource)：🌈 GitHub开源项目维护协同指南。
> - [read-list](https://github.com/eryajf/read-list)：📖 优质内容订阅，阅读方为根本
> - [awesome-github-profile-readme-chinese](https://github.com/eryajf/awesome-github-profile-readme-chinese)：🦩 优秀的中文区个人主页搜集

🚜 我还创建了一个项目 **[awesome-chatgpt-answer](https://github.com/eryajf/awesome-chatgpt-answer)** ：记录那些问得好，答得妙的时刻，欢迎提交你与ChatGPT交互过程中遇到的那些精妙对话。

⚗️ openai官方提供了一个 **[状态页](https://status.openai.com/)** 来呈现当前openAI服务的状态，同时如果有问题发布公告也会在这个页面，如果你感觉它有问题了，可以在这个页面看看。

## 功能介绍

- 🚀 帮助菜单：通过发送 `帮助` 将看到帮助列表，[🖼 查看示例](#%E5%B8%AE%E5%8A%A9%E5%88%97%E8%A1%A8)
- 🥷 私聊：支持与机器人单独私聊(无需艾特)，[🖼 查看示例](#%E4%B8%8E%E6%9C%BA%E5%99%A8%E4%BA%BA%E7%A7%81%E8%81%8A)
- 💬 群聊：支持在群里艾特机器人进行对话
- 🙋 单聊模式：每次对话都是一次新的对话，没有历史聊天上下文联系
- 🗣 串聊模式：带上下文理解的对话模式
- 🎨 图片生成：通过发送 `#图片`关键字开头的内容进行生成图片，[🖼 查看示例](#%E7%94%9F%E6%88%90%E5%9B%BE%E7%89%87)
- 🎭 角色扮演：支持场景模式，通过 `#周报` 的方式触发内置prompt模板 [🖼 查看示例](#%E9%80%9A%E8%BF%87%E5%86%85%E7%BD%AEprompt%E8%81%8A%E5%A4%A9)
- 🧑‍💻 频率限制：通过配置指定，自定义单个用户单日最大对话次数
- 💵 余额查询：通过发送 `余额` 关键字查询当前key所剩额度，[🖼 查看示例](#%E6%9F%A5%E8%AF%A2%E4%BD%99%E9%A2%9D)
- 🔗 自定义api域名：通过配置指定，解决国内服务器无法直接访问openai的问题
- 🪜 添加代理：通过配置指定，通过给应用注入代理解决国内服务器无法访问的问题
- 👐 默认模式：支持自定义默认的聊天模式，通过配置化指定
- 📝 查询对话：通过发送`#查对话 username:xxx`查询xxx的对话历史，可在线预览，可下载到本地
- 👹 白名单机制：通过配置指定，支持指定群组名称和用户名称作为白名单，从而实现可控范围与机器人对话
- 💂‍♀️ 管理员机制：通过配置指定管理员，部分敏感操作，以及一些应用配置，管理员有权限进行操作

## 使用前提

* 有Openai账号，并且创建好`api_key`，注册相关事项可以参考[此文章](https://juejin.cn/post/7173447848292253704) 。访问[这里](https://beta.openai.com/account/api-keys)，申请个人秘钥。
* 在钉钉开发者后台创建机器人，配置应用程序回调。

## 使用教程

### 第一步，创建机器人

#### 方案一：outgoing类型机器人

钉钉群内的机器人有一个outgoing模式，当你创建机器人的时候，可以选择启用这个模式，然后直接配置回调地址，免去在管理后台创建应用的步骤，就可以直接投入使用。

官方文档：[自定义机器人接入](https://open.dingtalk.com/document/orgapp/custom-robot-access)

但是这个模式貌似是部分开放的(目前来看貌似是部分人有创建这个类型的白名单)，所以如果你在钉钉群聊中添加`自定义机器人`的时候，看到和我一样的信息，则说明无法使用这种方式：

<details>
    <summary>🖼 点我查看示例图</summary>
    <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230325_162017.jpg">
</details>

`📢 注意`

- 如果你的和我一样，那么就只能放弃这种方案，往下看第二种对接方案。
- 如果使用这种方案，那么就不能与机器人私聊对话，只能局限在群聊当中艾特机器人聊天。
- 如果使用这种方案，则在群聊当中并不能达到真正的艾特发消息人的效果，因为这种机器人回调过来的关键信息为空。

#### 方案二：企业内部应用

创建步骤参考文档：[企业内部开发机器人](https://open.dingtalk.com/document/robots/enterprise-created-chatbot)，或者根据如下步骤进行配置。

1. 创建机器人。
    <details>
      <summary>🖼 点我查看示例图</summary>
      <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20221209_163616.png">
    </details>

   > `📢 注意1：`可能现在创建机器人的时候名字为`chatgpt`会被钉钉限制，请用其他名字命名。
   > `📢 注意2：`第四步骤点击创建应用的时候，务必选择使用旧版，从而创建旧版机器人。

   步骤比较简单，这里就不赘述了。

2. 配置机器人回调接口。
    <details>
      <summary>🖼 点我查看示例图</summary>
      <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20221209_163652.png">
    </details>

   创建完毕之后，点击机器人开发管理，然后配置将要部署的服务所在服务器的出口IP，以及将要给服务配置的域名。

  ` 如果提示：` 消息接收地址校验失败（请确保公网可访问该地址，如无有效SSL证书，可选择禁用证书校验），那么可以先输入一个`https://`，然后就能看到`禁用https`的选项了，选择禁用，然后再把地址改成`http`就好了。

3. 发布机器人。
    <details>
      <summary>🖼 点我查看示例图</summary>
      <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20221209_163709.png">
    </details>

   点击版本管理与发布，然后点击上线，这个时候就能在钉钉的群里中添加这个机器人了。

4. 群聊添加机器人。
    <details>
      <summary>🖼 点我查看示例图</summary>
      <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20221209_163724.png">
    </details>

### 第二步，部署应用

#### docker部署

你可以使用docker快速运行本项目。

```
第一种：基于环境变量运行
# 运行项目
$ docker run -itd --name chatgpt -p 8090:8090 \
  -v ./data:/app/data --add-host="host.docker.internal:host-gateway" \
  -e LOG_LEVEL="info" -e APIKEY=换成你的key -e BASE_URL="" \
  -e MODEL="gpt-3.5-turbo" -e SESSION_TIMEOUT=600 \
  -e HTTP_PROXY="http://host.docker.internal:15732" \
  -e DEFAULT_MODE="单聊" -e MAX_REQUEST=0 -e PORT=8090 \
  -e SERVICE_URL="你当前服务外网可访问的URL" -e CHAT_TYPE="0" \
  -e ALLOW_GROUPS=a,b -e ALLOW_USERS=a,b -e ADMIN_USERS=a,b -e APP_SECRETS="xxx,yyy" \
  -e AZURE_ON="false" -e AZURE_API_VERSION="" -e AZURE_RESOURCE_NAME="" \
  -e AZURE_DEPLOYMENT_NAME="" -e AZURE_OPENAI_TOKEN="" \
  -e HELP="欢迎使用本工具\n\n你可以查看：[用户指南](https://github.com/eryajf/chatgpt-dingtalk/blob/main/docs/userGuide.md)\n\n这是一个[开源项目](https://github.com/eryajf/chatgpt-dingtalk/)
  ，觉得不错你可以来波素质三连."  \
  --restart=always  dockerproxy.com/eryajf/chatgpt-dingtalk:latest
```

> 运行命令中映射的配置文件参考下边的[配置文件说明](#%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E8%AF%B4%E6%98%8E)。

- `📢 注意：`如果使用docker部署，那么PORT参数不需要进行任何调整。
- `📢 注意：`ALLOW_GROUPS,ALLOW_USERS,ADMIN_USERS三个参数为数组，如果需要指定多个，可用英文逗号分割。
- `📢 注意：`如果服务器节点本身就在国外或者自定义了`BASE_URL`，那么就把`HTTP_PROXY`参数留空即可。
- `📢 注意：`如果使用docker部署，那么proxy地址可以直接使用如上方式部署，`host.docker.internal`会指向容器所在宿主机的IP，只需要更改端口为你的代理端口即可。参见：[Docker容器如何优雅地访问宿主机网络](https://wiki.eryajf.net/pages/674f53/)


```
第二种：基于配置文件挂载运行
# 复制配置文件，根据自己实际情况，调整配置里的内容
$ cp config.example.yml config.yml  # 其中 config.example.yml 从项目的根目录获取

# 运行项目
$ docker run -itd --name chatgpt -p 8090:8090  -v `pwd`/config.yml:/app/config.yml --restart=always  dockerproxy.com/eryajf/chatgpt-dingtalk:latest
```

其中配置文件参考下边的配置文件说明。

```
第三种：使用 docker compose 运行
$ wget https://raw.githubusercontent.com/eryajf/chatgpt-dingtalk/main/docker-compose.yml

$ nano docker-compose.yml # 编辑 APIKEY 等信息

$ docker compose up -d
```

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

`📢 注意`:Nginx代理步骤是个可选步骤，你也可以直接通过服务器外网IP:PORT作为回调地址。

Nginx配置完毕之后，可以先手动请求一下，通过服务日志输出判断服务是否正常可用：

```sh
$ curl --location --request POST 'http://chat.eryajf.net/' \
  --header 'Content-type: application/json' \
  -d '{
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

> **📢 注意：** 如果配置文件中`app_secret`不为空，那么这里curl请求将会校验失败，理论上只要服务能够正常被访问，那么直接在钉钉管理后台配置回调就可以了，如果想通过curl进行调试，则注意需把对应配置项留空。

#### 二进制部署


如果你想通过命令行直接部署，可以直接下载release中的[压缩包](https://github.com/eryajf/chatgpt-dingtalk/releases) ，请根据自己系统以及架构选择合适的压缩包，下载之后直接解压运行。

下载之后，在本地解压，即可看到可执行程序，与配置文件：

```sh
$ tar xf chatgpt-dingtalk-v0.0.4-darwin-arm64.tar.gz
$ cd chatgpt-dingtalk-v0.0.4-darwin-arm64
$ cp config.example.yml  config.yml # 然后根据情况调整配置文件内容,宿主机如遇端口冲突,可通过调整config.yml中的port参数自定义服务端口
$ ./chatgpt-dingtalk  # 直接运行

# 如果要守护在后台运行
$ nohup ./chatgpt-dingtalk &> run.log &
$ tail -f run.log
```

## 亮点特色

### 与机器人私聊

`2023-03-08`补充，我发现也可以不在群里艾特机器人聊天，还可点击机器人，然后点击发消息，通过与机器人直接对话进行聊天：

> 由 [@Raytow](https://github.com/Raytow) 同学发现，在机器人自动生成的测试群里无法直接私聊机器人，在其他群里单独添加这个机器人，然后再点击就可以跟它私聊了。

<details>
    <summary>🖼 点我查看示例图</summary>
    <img src="https://user-images.githubusercontent.com/33259379/223607306-2ac836a2-7ce5-4a12-a16e-bec40b22d8d6.png">
</details>


### 帮助列表

> 艾特机器人发送空内容或者帮助，会返回帮助列表。

<details>
    <summary>🖼 点我查看示例图</summary>
    <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230216_221253.png">
</details>

### 切换模式

> 发送指定关键字，可以切换不同的模式。

<details>
    <summary>🖼 点我查看示例图</summary>
    <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230215_184655.png">
</details>

> 📢 注意：串聊模式下，群里每个人的聊天上下文是独立的。
> 📢 注意：默认对话模式为单聊，因此不必发送单聊即可进入单聊模式，而要进入串聊，则需要发送串聊关键字进行切换，当串聊内容超过最大限制的时候，你可以发送重置，然后再次进入串聊模式。

### 查询余额

> 艾特机器人发送 `余额` 二字，会返回当前key对应的账号的剩余额度以及可用日期。

<details>
    <summary>🖼 点我查看示例图</summary>
    <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230304_222522.jpg">
</details>

### 日常问题

<details>
    <summary>🖼 点我查看示例图</summary>
    <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20221209_163739.png">
</details>

### 通过内置prompt聊天

> 发送模板两个字，会返回当前内置支持的prompt列表。

<details>
    <summary>🖼 点我查看示例图</summary>
    <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230323_152703.jpg">
</details>

> 如果你发现有比较优秀的prompt，欢迎PR。注意：一些与钉钉使用场景不是很匹配的，就不要提交了。

### 生成图片

> 发送以 `#图片`开头的内容，将会触发绘画能力，图片生成之后，将会保存在程序根目录下的`images目录`下。
>
> 如果你绘图没有思路，可以在[这里](https://www.clickprompt.org/zh-CN/)以及[这里](https://lexica.art/)找到一些不错的prompt。

<details>
    <summary>🖼 点我查看示例图</summary>
    <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230323_150547.jpg">
</details>

### 支持 gpt-4

如果你的账号通过了官方的白名单，那么可以将模型配置为：`gpt-4-0314`或`gpt-4`，目前gpt-4的余额查询以及图片生成功能暂不可用，可能是接口限制，也可能是其他原因，等我有条件的时候，会对这些功能进行测试验证。

> 以下是gpt-3.5与gpt-4对数学计算方面的区别。

<details>
    <summary>🖼 点我查看示例图</summary>
    <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230330_180308.jpg">
</details>

感谢[@PIRANHACHAN](https://github.com/PIRANHACHAN)同学提供的gpt-4的key，使得项目在gpt-4的对接上能够进行验证测试，达到了可用状态。

##  本地开发

```sh
# 获取项目
$ git clone https://github.com/eryajf/chatgpt-dingtalk.git

# 进入项目目录
$ cd chatgpt-dingtalk

# 复制配置文件，根据个人实际情况进行配置
$ cp config.example.yml config.yml

# 启动项目
$ go run main.go
```

## 配置文件说明

```yaml
# 应用的日志级别，info or debug
log_level: "info"
# openai api_key
api_key: "xxxxxxxxx"
# 如果你使用官方的接口地址 https://api.openai.com，则留空即可，如果你想指定请求url的地址，可通过这个参数进行配置，注意需要带上 http 协议
base_url: ""
# 指定模型，默认为 gpt-3.5-turbo , 可选参数有： "gpt-4-0314", "gpt-4", "gpt-3.5-turbo-0301", "gpt-3.5-turbo"，如果使用gpt-4，请确认自己是否有接口调用白名单
model: "gpt-3.5-turbo"
# 会话超时时间,默认600秒,在会话时间内所有发送给机器人的信息会作为上下文
session_timeout: 600
# 指定请求时使用的代理，如果为空，则不使用代理，注意需要带上 http 协议 或 socks5 协议
http_proxy: ""
# 指定默认的对话模式，可根据实际需求进行自定义，如果不设置，默认为单聊，即无上下文关联的对话模式
default_mode: "单聊"
# 单人单日请求次数上限，默认为0，即不限制
max_request: 0
# 指定服务启动端口，默认为 8090，一般在二进制宿主机部署时，遇到端口冲突时使用
port: "8090"
# 指定服务的地址，就是当前服务可供外网访问的地址(或者直接理解为你配置在钉钉回调那里的地址)，用于生成图片时给钉钉做渲染
service_url: "http://chat.eryajf.net"
# 限定对话类型 0：不限 1：只能单聊 2：只能群聊
chat_type: "0"
# 哪些群组可以进行对话，如果留空，则表示允许所有群组，如果要限制，则写群组的名称，比如 ["aa","bb"]
# 对话聊天时，如下三个满足其一即可通过校验
allow_groups:
  - "学无止境"
# 哪些用户可以进行对话，如果留空，则表示允许所有用户，如果要限制，则列表中写用户的名称，比如 ["张三","李四"]
allow_users: ["张三","李四"]
# 指定哪些人为此系统的管理员，如果留空，则表示没有人是管理员，如果要限制，则列表中写用户的userid
# 比如 ["1301691029702722","1301691029702733"]，这个信息需要在钉钉管理后台的通讯录当中获取：https://oa.dingtalk.com/contacts.htm#/contacts
# 注意：如果下边的app_secrets为空，以及使用outgoing的方式配置机器人，这两种情况下，都表示没有人是管理员
admin_users: []
# 钉钉机器人在应用信息中的AppSecret，为了校验回调的请求是否合法，如果留空，将会忽略校验，则该接口将会存在其他人也能随意调用的安全隐患，因此强烈建议配置正确的secret，如果你的服务对接给多个机器人，这里可以配置多个机器人的secret
app_secrets: []
# 帮助信息，放在配置文件，可供自定义
help: "欢迎使用本工具\n\n你可以查看：[用户指南](https://github.com/eryajf/chatgpt-dingtalk/blob/main/docs/userGuide.md)\n\n这是一个[开源项目](https://github.com/eryajf/chatgpt-dingtalk/)，觉得不错你可以来波素质三连."

# Azure OpenAI 配置
# 例如 curl https://forkway-test.openai.azure.com/openai/deployments/test-dev/chat/completions\?api-version\=2023-03-15-preview \
# azure_api_version: "2023-03-15-preview"
# azure_resource_name: "forkway-test"
# azure_deployment_name: "test-dev"
azure_on: false # 如果是true，则会走azure的openai接口
azure_api_version: "2023-03-15-preview"
azure_resource_name: "xxxx"
azure_deployment_name: "xxxx"
azure_openai_token: "xxxx"

```

## 常见问题

如何更好地使用ChatGPT：这里有[许多案例](https://github.com/f/awesome-chatgpt-prompts)可供参考。

`🗣 重要重要` 一些常见的问题，我单独开issue放在这里：[👉点我👈](https://github.com/eryajf/chatgpt-dingtalk/issues/44)，可以查看这里辅助你解决问题，如果里边没有，请对历史issue进行搜索(不要提交重复的issue)，也欢迎大家补充。

## 进群交流

我创建了一个钉钉的交流群，欢迎进群交流。

![](https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230405_191425.jpg)

## 感谢

这个项目能够成立，离不开这些开源项目：

- [go-resty/resty](https://github.com/go-resty/resty)
- [patrickmn/go-cache](https://github.com/patrickmn/go-cache)
- [solywsh/chatgpt](https://github.com/solywsh/chatgpt)
- [xgfone/ship](https://github.com/xgfone/ship)
- [avast/retry-go](https://github.com/avast/retry-go)
- [sashabaranov/go-openapi](https://github.com/sashabaranov/go-openai)
- [charmbracelet/log](https://github.com/charmbracelet/log)

## 赞赏

如果觉得这个项目对你有帮助，你可以请作者[喝杯咖啡 ☕️](https://wiki.eryajf.net/reward/)

## 高光时刻

> 本项目曾在 | [2022-12-12](https://github.com/bonfy/github-trending/blob/master/2022/2022-12-12.md#go) | [2022-12-18](https://github.com/bonfy/github-trending/blob/master/2022/2022-12-18.md#go) | [2022-12-19](https://github.com/bonfy/github-trending/blob/master/2022/2022-12-19.md#go) | [2022-12-20](https://github.com/bonfy/github-trending/blob/master/2022/2022-12-20.md#go) | [2023-02-09](https://github.com/bonfy/github-trending/blob/master/2023-02-09.md#go) | [2023-02-10](https://github.com/bonfy/github-trending/blob/master/2023-02-10.md#go) | [2023-02-11](https://github.com/bonfy/github-trending/blob/master/2023-02-11.md#go) | [2023-02-12](https://github.com/bonfy/github-trending/blob/master/2023-02-12.md#go) | [2023-02-13](https://github.com/bonfy/github-trending/blob/master/2023-02-13.md#go) | [2023-02-14](https://github.com/bonfy/github-trending/blob/master/2023-02-14.md#go) | [2023-02-15](https://github.com/bonfy/github-trending/blob/master/2023-02-15.md#go) | [2023-03-04](https://github.com/bonfy/github-trending/blob/master/2023-03-04.md#go) | [2023-03-05](https://github.com/bonfy/github-trending/blob/master/2023-03-05.md#go) | [2023-03-19](https://github.com/bonfy/github-trending/blob/master/2023-03-19.md#go) | [2023-03-22](https://github.com/bonfy/github-trending/blob/master/2023-03-22.md#go) | [2023-03-25](https://github.com/bonfy/github-trending/blob/master/2023-03-25.md#go) | [2023-03-26](https://github.com/bonfy/github-trending/blob/master/2023-03-26.md#go) | [2023-03-27](https://github.com/bonfy/github-trending/blob/master/2023-03-27.md#go) | [2023-03-29](https://github.com/bonfy/github-trending/blob/master/2023-03-29.md#go), 这些天里，登上GitHub Trending。而且还在持续登榜中，可见最近openai的热度。
> ![image_20230316_114915](https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230316_114915.jpg)

## 贡献者列表

<!-- readme: collaborators,contributors -start -->
<table>
<tr>
    <td align="center">
        <a href="https://github.com/eryajf">
            <img src="https://avatars.githubusercontent.com/u/33259379?v=4" width="100;" alt="eryajf"/>
            <br />
            <sub><b>二丫讲梵</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/ffinly">
            <img src="https://avatars.githubusercontent.com/u/29793346?v=4" width="100;" alt="ffinly"/>
            <br />
            <sub><b>Finly</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/Leizhenpeng">
            <img src="https://avatars.githubusercontent.com/u/50035229?v=4" width="100;" alt="Leizhenpeng"/>
            <br />
            <sub><b>RiverRay</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/laorange">
            <img src="https://avatars.githubusercontent.com/u/68316902?v=4" width="100;" alt="laorange"/>
            <br />
            <sub><b>辣橙</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/luoxufeiyan">
            <img src="https://avatars.githubusercontent.com/u/6621172?v=4" width="100;" alt="luoxufeiyan"/>
            <br />
            <sub><b>Hugh Gao</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/fantasticmao">
            <img src="https://avatars.githubusercontent.com/u/20675747?v=4" width="100;" alt="fantasticmao"/>
            <br />
            <sub><b>Mao Mao</b></sub>
        </a>
    </td></tr>
<tr>
    <td align="center">
        <a href="https://github.com/iblogc">
            <img src="https://avatars.githubusercontent.com/u/3283023?v=4" width="100;" alt="iblogc"/>
            <br />
            <sub><b>Iblogc</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/WinMin">
            <img src="https://avatars.githubusercontent.com/u/18380453?v=4" width="100;" alt="WinMin"/>
            <br />
            <sub><b>Swing</b></sub>
        </a>
    </td></tr>
</table>
<!-- readme: collaborators,contributors -end -->
