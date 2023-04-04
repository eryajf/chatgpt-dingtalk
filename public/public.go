package public

import (
	"github.com/eryajf/chatgpt-dingtalk/config"
	"github.com/eryajf/chatgpt-dingtalk/pkg/cache"
	"github.com/eryajf/chatgpt-dingtalk/pkg/db"
)

var UserService cache.UserServiceInterface
var Config *config.Configuration
var Prompt *[]config.Prompt

func InitSvc() {
	// 加载配置
	Config = config.LoadConfig()
	// 加载prompt
	Prompt = config.LoadPrompt()
	// 初始化缓存
	UserService = cache.NewUserService()
	// 初始化数据库
	db.InitDB()
	// 暂时不在初始化时获取余额
	// if Config.Model == openai.GPT3Dot5Turbo0301 || Config.Model == openai.GPT3Dot5Turbo {
	// _, _ = GetBalance()
	// }
}

var Welcome string = `# 发送信息

若您想给机器人发送信息，有如下两种方式：

1. 在本机器人所在群里@机器人；
2. 点击机器人的头像后，再点击"发消息"。

机器人收到您的信息后，默认会交给chatgpt进行处理。除非，您发送的内容是如下**系统指令**之一。

-----

# 系统指令

系统指令是一些特殊的词语，当您向机器人发送这些词语时，会触发对应的功能：

**单聊**：每条消息都是单独的对话，不包含上下文

**串聊**：对话会携带上下文，除非您主动重置对话或对话长度超过限制

**重置**：重置上下文

**余额**： ~~查询机器人所用OpenAI账号的余额~~ (暂不可用)

**模板**：查询机器人内置的快捷模板

**图片**：查看如何根据提示词生成图片

**查对话**：获取指定人员的对话历史

**帮助**：重新获取帮助信息

-----

# 友情提示

使用"串聊模式"会显著加快机器人所用账号的余额消耗速度。

因此，若无保留上下文的需求，建议使用"单聊模式"。

即使有保留上下文的需求，也应适时使用"重置"指令来重置上下文。

-----

# 项目地址

本项目已在GitHub开源，[查看源代码](https://github.com/eryajf/chatgpt-dingtalk)。
`
