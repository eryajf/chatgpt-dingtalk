# 流式输出功能使用指南

## 功能概述

项目已支持钉钉机器人的流式输出功能，可以让 AI 回答像打字一样逐字显示，提供更好的用户体验。

## 两种流式模式

### 1. 简化流式模式（推荐快速开始）

**特点**:
- 无需配置钉钉卡片模板
- 直接使用累积内容一次性回复
- 配置简单，开箱即用

**配置方式**:

在 `config.yml` 中添加:

```yaml
# 启用流式输出
stream_mode: true
```

### 2. 高级流式卡片模式

**特点**:
- 使用钉钉互动卡片实现真正的流式更新
- 内容逐步显示，类似 ChatGPT 网页版效果
- 需要在钉钉开放平台创建卡片模板

**配置方式**:

在 `config.yml` 中添加:

```yaml
# 启用流式输出
stream_mode: true
# 钉钉卡片模板ID (可选，用于高级流式卡片模式)
card_template_id: "your-card-template-id"
```

## 配置钉钉卡片模板（高级模式）

### 步骤 1: 创建卡片模板

1. 登录 [钉钉开放平台](https://open.dingtalk.com/)
2. 进入你的应用 -> 互动卡片 -> 卡片模板管理
3. 创建新模板，使用以下 JSON Schema:

```json
{
  "type": "object",
  "properties": {
    "content": {
      "type": "string",
      "title": "内容"
    }
  }
}
```

### 步骤 2: 设计卡片样式

在卡片编辑器中，添加一个 Markdown 组件来显示 `content` 字段:

```json
{
  "type": "markdown",
  "text": "{{content}}"
}
```

### 步骤 3: 发布并获取模板ID

1. 保存并发布卡片模板
2. 复制模板ID（类似: `4d18414c-aabc-4ec8-9e67-4ceefeada72a.schema`）
3. 将模板ID填入 `config.yml` 的 `card_template_id` 字段

## 完整配置示例

```yaml
# 日志级别
log_level: "info"

# OpenAI 配置
api_key: "sk-..."
model: "gpt-4o"
base_url: ""  # 可选，用于 API 中转

# 流式输出配置
stream_mode: true                           # 启用流式输出
card_template_id: ""                        # 可选：钉钉卡片模板ID

# 其他配置...
session_timeout: 600s
max_question_len: 2048
max_answer_len: 2048
max_text: 4096
default_mode: "单聊"
```

## 实现原理

### 简化模式流程

```
用户提问 → OpenAI 流式响应 → 累积完整内容 → 一次性回复
```

### 高级卡片模式流程

```
用户提问
  ↓
创建钉钉卡片（空内容）
  ↓
发送初始状态 "稍等，让我想一想..."
  ↓
OpenAI 流式响应
  ↓
接收到内容 → 立即累积到缓冲区
  ↓
距离上次更新超过300ms? → 是 → 更新卡片
  ↓ 否               ↓
继续接收 ←-----------┘
  ↓
流式结束，发送最终内容（标记为完成）
```

## 技术架构

### 新增文件

1. **[pkg/llm/stream.go](pkg/llm/stream.go)**
   - 实现 OpenAI 流式响应
   - 提供 `SingleQaStream()` 和 `ContextQaStream()` API
   - 支持 ChatCompletion 流式调用

2. **[pkg/dingbot/stream.go](pkg/dingbot/stream.go)**
   - 实现钉钉流式卡片更新
   - 封装钉钉 Streaming Update API
   - 提供 `UpdateAIStreamCard()` 方法

3. **[pkg/process/stream.go](pkg/process/stream.go)**
   - 实现流式处理逻辑
   - `DoStream()` - 简化流式模式
   - `DoStreamWithCard()` - 高级卡片模式
   - 包含定时更新和错误处理

### 改动文件

1. **[config/config.go](config/config.go)**
   - 添加 `StreamMode` 配置项
   - 添加 `CardTemplateID` 配置项

2. **[pkg/process/process_request.go](pkg/process/process_request.go)**
   - 根据配置自动选择流式或普通模式
   - 支持流式卡片和流式普通两种方式

## API 使用示例

### 在代码中使用流式 API

```go
import "github.com/eryajf/chatgpt-dingtalk/pkg/llm"

// 单聊流式
contentCh, cleanup, err := llm.SingleQaStream("你好", "user123")
if err != nil {
    log.Fatal(err)
}
defer cleanup()

for content := range contentCh {
    fmt.Print(content)  // 逐块输出
}

// 串聊流式
client, contentCh, err := llm.ContextQaStream("继续", "user123")
if err != nil {
    log.Fatal(err)
}
defer client.Close()

fullAnswer := ""
for content := range contentCh {
    fullAnswer += content
    fmt.Print(content)
}

// 保存对话上下文
client.ChatContext.SaveConversation("user123")
```

## 性能优化

### 流式更新策略

高级卡片模式采用**实时流式更新**策略:

- 从大模型接收到内容后立即更新卡片
- 使用缓冲机制避免更新过于频繁(默认最小间隔 **300ms**)
- 这样可以实现真正的实时流式体验,类似 ChatGPT 网页版

可以在 [pkg/process/stream.go](pkg/process/stream.go) 中修改最小更新间隔:

```go
minUpdateInterval := 300 * time.Millisecond  // 修改这里
```

建议范围：200ms - 500ms
- 更小的间隔:更实时,但 API 调用更频繁
- 更大的间隔:API 调用较少,但流式感觉不明显

### 流式模式选择建议

| 场景 | 推荐模式 | 原因 |
|------|---------|------|
| 快速部署 | 简化模式 | 无需额外配置 |
| 追求体验 | 高级卡片模式 | 真正的流式显示 |
| 高频使用 | 简化模式 | 减少 API 调用 |
| 演示展示 | 高级卡片模式 | 视觉效果更好 |

## 兼容性

- ✅ 保持原有非流式模式完全兼容
- ✅ 支持单聊和串聊两种对话模式
- ✅ 支持所有 OpenAI 兼容的模型
- ✅ 支持 Azure OpenAI
- ✅ 保留敏感词过滤、请求限制等功能

## 故障排查

### 问题 1: 流式模式不生效

**检查项**:
1. 确认 `config.yml` 中 `stream_mode: true`
2. 重启应用以加载新配置
3. 查看日志是否有错误信息

### 问题 2: 卡片模式无法显示 / 日志显示 "robot code is empty"

**原因**: 这是正常的降级行为

**说明**:
- 高级卡片模式需要通过 `credentials` 配置才能工作
- 如果没有配置 `credentials`，系统会自动降级为简化流式模式
- 简化流式模式不需要卡片，依然可以正常工作

**解决方案**:
1. 如果想使用高级卡片模式，需要在 `config.yml` 中配置 `credentials`:
   ```yaml
   credentials:
     - client_id: "your-app-key"
       client_secret: "your-app-secret"
   ```
2. 如果不需要卡片模式，可以忽略这个警告，或者将 `card_template_id` 留空

### 问题 3: 卡片模板配置正确但不显示

**检查项**:
1. 确认卡片模板ID正确
2. 确认卡片模板已发布
3. 确认钉钉应用有卡片权限
4. 确认 `credentials` 配置正确
5. 查看日志是否有降级提示

### 问题 4: 流式响应中断

**可能原因**:
1. OpenAI API 超时 - 检查网络连接
2. 钉钉 Access Token 过期 - 会自动刷新
3. 上下文超过限制 - 减少 `max_text` 配置

### 问题 5: HTTP/2 流式错误 "stream error: INTERNAL_ERROR"

**错误信息**:
```
Post "https://api.xxx.com/v1/completions": stream error: stream ID 5; INTERNAL_ERROR; received from peer
```

**可能原因**:
1. 上游 API 服务器内部错误或资源不足
2. 网络不稳定,长连接中断
3. HTTP/2 连接管理问题

**解决方案**:
1. **已优化**: 代码已优化 HTTP 客户端配置,增加连接池和超时设置
2. **部分内容保护**: 如果已接收到部分内容,不会因错误而丢失
3. **禁用 HTTP/2** (如果问题频繁): 在 [pkg/llm/client.go](pkg/llm/client.go) 中取消注释:
   ```go
   Transport: &http.Transport{
       // ...
       ForceAttemptHTTP2: false,  // 取消注释这行
   }
   ```
   注意:禁用 HTTP/2 会降低性能,但可能更稳定

4. **检查 API 服务器**: 如果使用中转服务,检查中转服务器的健康状况和资源

## 智能降级机制

系统实现了智能降级机制，确保即使高级功能无法使用，基础功能依然可用:

```
尝试高级卡片模式
  ↓
检查 RobotCode 是否存在
  ↓ (否)
降级为简化流式模式
  ↓
检查 credentials 配置
  ↓ (无配置)
降级为简化流式模式
  ↓
尝试创建卡片
  ↓ (失败)
降级为简化流式模式
  ↓
正常流式输出
```

这意味着：
- ✅ 即使配置不完整，流式功能依然可用
- ✅ 不会因为卡片失败而导致整个功能不可用
- ✅ 日志会清楚地告诉你当前使用的是哪种模式

## 参考资料

- [钉钉流式消息更新 API](https://open.dingtalk.com/document/development/api-streamingupdate)
- [OpenAI Streaming API](https://platform.openai.com/docs/api-reference/streaming)
- [PandaWiki 项目参考实现](tmp/PandaWiki)

## 贡献

欢迎提交 Issue 和 Pull Request 来改进流式功能！
