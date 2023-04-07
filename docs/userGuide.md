本文是chatgpt-dingtalk项目的使用指南，该项目涉及的指令，以及特性，都会在本文呈现。

## 发送信息

若您想给机器人发送信息，有如下两种方式：

1. **群聊：** 在机器人所在群里`@机器人` 后边跟着要提问的内容。
2. **私聊：** 点击机器人的`头像`后，再点击`发消息`。

## 系统指令

系统指令是一些特殊的词语，当您向机器人发送这些词语时，会触发对应的功能。

**📢 注意：系统指令，即只发指令，没有特殊标识，也没有内容。**

以下是系统指令详情：


|    指令    |                     描述                     |                             示例                             | 补充 |
| :--------: | :------------------------------------------: | :----------------------------------------------------------: | :--: |
|  **单聊**  | 每次对话都是一次新的对话，没有聊天上下文联系 |    <details><br /><summary>点击查看</summary><br /><img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230404_193608.jpg"><br /></details>                                                          |      |
|  **串聊**  |            带上下文联系的对话模式            |     <details><br /><summary>点击查看</summary><br /><img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230404_193608.jpg"><br /></details>                                                         |      |
|  **重置**  |        重置上下文模式，回归到默认模式        |        <details><br /><summary>点击查看</summary><br /><img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230404_193608.jpg"><br /></details>                                                      |      |
|  **余额**  | 查询机器人所用OpenAI账号的余额 |       <details><br /><summary>点击查看</summary><br /><img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230304_222522.jpg"><br /></details>                                                       |      |
|  **模板**  |           查看应用内置的prompt模板           |      <details><br /><summary>点击查看</summary><br /><img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230404_193827.jpg"><br /></details>                                                        |      |
|  **图片**  |          查看如何根据提示生成图片          | <details><br /><summary>点击查看</summary><br /><img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230404_194125.jpg"><br /></details> |      |
| **查对话** |            获取指定人员的对话历史            |      <details><br /><summary>点击查看</summary><br /><img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230404_193938.jpg"><br /></details>                                                        |      |
|  **帮助**  |                 获取帮助信息                 |     <details><br /><summary>点击查看</summary><br /><img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230404_202336.jpg"><br /></details>                                                         |      |

## 功能指令

除去系统指令，还有一些功能指令，功能指令是直接与应用交互，达到交互目的的一种指令。

**📢 注意：功能指令，一律以 #+关键字 为开头，通常需要在关键字后边加个空格，然后再写描述或参数。**

以下是功能指令详情

| 指令 | 说明 | 示例 | 补充|
| :--: | :--: | :--: | :--: |
|  **#图片**  |          根据提示咒语生成对应图片          | <details><br /><summary>点击查看</summary><br /><img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230323_150547.jpg"><br /></details> |      |
| **#域名**     | 查询域名相关信息     |  <details><br /><summary>点击查看</summary><br /><img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230404_202620.jpg"><br /></details>    |  |
| **#证书**     | 查询域名证书相关信息     | <details><br /><summary>点击查看</summary><br /><img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230404_202706.jpg"><br /></details>    |  |
| **#Linux命令**     | 根据自然语言描述生成对应命令     | <details><br /><summary>点击查看</summary><br /><img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230404_214947.jpg"><br /></details>    | 此指令中的Linux开头字幕可以大写 |
| **#解释代码**     | 分析一段代码的功能或含义     | <details><br /><summary>点击查看</summary><br /><img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230404_215242.jpg"><br /></details>    |  |
| **#正则**     | 根据自然语言描述生成正则     | <details><br /><summary>点击查看</summary><br /><img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230404_220222.jpg"><br /></details>    |  |
| **#周报**     | 应用周报的prompt     | <details><br /><summary>点击查看</summary><br /><img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230404_214335.jpg"><br /></details>    |  |
| **#生成sql**     | 根据自然语言描述生成sql语句     | <details><br /><summary>点击查看</summary><br /><img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230404_221325.jpg"><br /></details>    |  |

如上大多数能力，都是依赖prompt模板实现，如果你有更好的prompt，欢迎提交PR。

## 友情提示

使用`串聊模式`会显著加快机器人所用账号的余额消耗速度，因此，若无保留上下文的需求，建议使用`单聊模式`。

即使有保留上下文的需求，也应适时使用`重置`指令来重置上下文。

## 项目地址

本项目已在GitHub开源，[查看源代码](https://github.com/eryajf/chatgpt-dingtalk)。