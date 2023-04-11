# 贡献者指南

欢迎反馈、bug报告和拉取请求，可点击[issue](https://github.com/eryajf/chatgpt-dingtalk/issues) 提交.

如果你是第一次进行GitHub的协作，可参阅： [协同开发流程](https://eryajf.github.io/HowToStartOpenSource/views/01-basic-content/03-collaborative-development-process.html)

## 注意事项

- 如果你的变更中新增或者减少了配置，那么需要注意有这么几个地方需要同步调整：

  - [config.go](https://github.com/eryajf/chatgpt-dingtalk/blob/main/config/config.go)
  - [config.example.yml](https://github.com/eryajf/chatgpt-dingtalk/blob/main/config.example.yml)
  - [docker-compose.yml](https://github.com/eryajf/chatgpt-dingtalk/blob/main/docker-compose.yml)
  - [README.md](https://github.com/eryajf/chatgpt-dingtalk/blob/main/README.md)
    - docker [启动命令](https://github.com/eryajf/chatgpt-dingtalk/blob/main/README.md#docker%E9%83%A8%E7%BD%B2)中要添加。
    - [配置文件说明](https://github.com/eryajf/chatgpt-dingtalk/blob/main/README.md#%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E8%AF%B4%E6%98%8E)要添加。

  一定要检查这几个地方，而且注意务必加好注释，否则将会影响用户升级体验新功能。

- 关于配置管理还有一个很重要的点在于，`config.example.yml`中的配置务必配置为最大权限，以免出现用户首次部署就无法走到正常逻辑的情况。

- 请务必检查你的提交，是否包含secret，api_key之类的信息，如果要贴示例，注意数据脱敏。

- 如果新增了功能性的模板，则务必在[使用指南](./docs/userGuide.md)中添加对应说明文档。