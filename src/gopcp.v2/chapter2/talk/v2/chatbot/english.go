package chatbot

import (
	"fmt"
	"strings"
)

// simpleEN 代表针对英文的演示级聊天机器人。
type simpleEN struct {
	name string
	talk Talk
}

// NewSimpleEN 用于创建针对英文的演示级聊天机器人。
func NewSimpleEN(name string, talk Talk) Chatbot {
	return &simpleEN{
		name: name,
		talk: talk,
	}
}

// Name 是Chatbot接口的实现的一部分。
func (robot *simpleEN) Name() string {
	return robot.name
}

// Begin 是Chatbot接口的实现的一部分。
func (robot *simpleEN) Begin() (string, error) {
	return "Please input your name:", nil
}

// Hello 是Talk接口的实现的一部分。
func (robot *simpleEN) Hello(userName string) string {
	userName = strings.TrimSpace(userName)
	if robot.talk != nil {
		return robot.talk.Hello(userName)
	}
	return fmt.Sprintf("Hello, %s! What can I do for you?", userName)
}

// Talk 是Talk接口的实现的一部分。
func (robot *simpleEN) Talk(heard string) (saying string, end bool, err error) {
	heard = strings.TrimSpace(heard)
	if robot.talk != nil {
		return robot.talk.Talk(heard)
	}
	switch heard {
	case "":
		return
	case "nothing", "bye":
		saying = "Bye!"
		end = true
		return
	default:
		saying = "Sorry, I didn't catch you."
		return
	}
}

// ReportError 是Chatbot接口的实现的一部分。
func (robot *simpleEN) ReportError(err error) string {
	return fmt.Sprintf("An error occurred: %s\n", err)
}

// End 是Chatbot接口的实现的一部分。
func (robot *simpleEN) End() error {
	return nil
}
