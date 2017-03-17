package chatbot

import "errors"

// Talk 定义了聊天的接口类型。
type Talk interface {
	Hello(userName string) string
	Talk(heard string) (saying string, end bool, err error)
}

// Chatbot 定义了聊天机器人的接口类型。
type Chatbot interface {
	Name() string
	Begin() (string, error)
	Talk
	ReportError(err error) string
	End() error
}

var (
	// ErrInvalidChatbotName 代表无效的聊天机器人名称的错误。
	ErrInvalidChatbotName = errors.New("Invalid chatbot name")
	// ErrInvalidChatbot 代表无效的聊天机器人的错误。
	ErrInvalidChatbot = errors.New("Invalid chatbot")
	// ErrExistingChatbot 代表对同名的聊天机器人进行重复注册的错误。
	ErrExistingChatbot = errors.New("Existing chatbot")
)

// chatbotMap 代表了名称-聊天机器人的映射。
var chatbotMap = map[string]Chatbot{}

// Register 用于注册聊天机器人。
// 若结果值为nil则表示注册成功。
func Register(chatbot Chatbot) error {
	if chatbot == nil {
		return ErrInvalidChatbot
	}
	name := chatbot.Name()
	if name == "" {
		return ErrInvalidChatbotName
	}
	if _, ok := chatbotMap[name]; ok {
		return ErrExistingChatbot
	}
	chatbotMap[name] = chatbot
	return nil
}

// Get 用于获取指定名称的聊天机器人。
func Get(name string) Chatbot {
	return chatbotMap[name]
}
