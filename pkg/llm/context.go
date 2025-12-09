package llm

import (
	"bytes"
	"encoding/gob"
	"strings"

	"github.com/eryajf/chatgpt-dingtalk/public"
)

var (
	DefaultAiRole    = "AI"
	DefaultHumanRole = "Human"

	DefaultCharacter  = []string{"helpful", "creative", "clever", "friendly", "lovely", "talkative"}
	DefaultBackground = "The following is a conversation with AI assistant. The assistant is %s"
	DefaultPreset     = "\n%s: 你好,让我们开始愉快的谈话!\n%s: 我是 AI assistant ,请问你有什么问题?"
)

type Context struct {
	background  string
	preset      string
	maxSeqTimes int
	aiRole      *role
	humanRole   *role

	old        []conversation
	restartSeq string
	startSeq   string

	seqTimes int

	maintainSeqTimes bool
}

type ContextOption func(*Context)

type conversation struct {
	Role   *role
	Prompt string
}

type role struct {
	Name string
}

func NewContext(options ...ContextOption) *Context {
	ctx := &Context{
		aiRole:           &role{Name: DefaultAiRole},
		humanRole:        &role{Name: DefaultHumanRole},
		background:       "",
		maxSeqTimes:      1000,
		preset:           "",
		old:              []conversation{},
		seqTimes:         0,
		restartSeq:       "\n" + DefaultHumanRole + ": ",
		startSeq:         "\n" + DefaultAiRole + ": ",
		maintainSeqTimes: false,
	}

	for _, option := range options {
		option(ctx)
	}
	return ctx
}

func (c *Context) PollConversation() {
	c.old = c.old[1:]
	c.seqTimes--
}

func (c *Context) ResetConversation(userid string) {
	public.UserService.ClearUserSessionContext(userid)
}

func (c *Context) SaveConversation(userid string) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(c.old)
	if err != nil {
		return err
	}
	public.UserService.SetUserSessionContext(userid, buffer.String())
	return nil
}

func (c *Context) LoadConversation(userid string) error {
	dec := gob.NewDecoder(strings.NewReader(public.UserService.GetUserSessionContext(userid)))
	err := dec.Decode(&c.old)
	if err != nil {
		return err
	}
	c.seqTimes = len(c.old)
	return nil
}

func (c *Context) SetHumanRole(role string) {
	c.humanRole.Name = role
	c.restartSeq = "\n" + c.humanRole.Name + ": "
}

func (c *Context) SetAiRole(role string) {
	c.aiRole.Name = role
	c.startSeq = "\n" + c.aiRole.Name + ": "
}

func (c *Context) SetMaxSeqTimes(times int) {
	c.maxSeqTimes = times
}

func (c *Context) GetMaxSeqTimes() int {
	return c.maxSeqTimes
}

func (c *Context) SetBackground(background string) {
	c.background = background
}

func (c *Context) SetPreset(preset string) {
	c.preset = preset
}

func WithMaxSeqTimes(times int) ContextOption {
	return func(c *Context) {
		c.SetMaxSeqTimes(times)
	}
}

func WithOldConversation(userid string) ContextOption {
	return func(c *Context) {
		_ = c.LoadConversation(userid)
	}
}

func WithMaintainSeqTimes(maintain bool) ContextOption {
	return func(c *Context) {
		c.maintainSeqTimes = maintain
	}
}
