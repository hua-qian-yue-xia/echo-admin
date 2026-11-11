package rep

// Rep 是统一的响应体
type Rep[T any] struct {
	RequestID string       `json:"requestId"`        // RequestID 与响应头 X-Request-ID 一致
	Code      string       `json:"code"`             // Code 是稳定的机器可读业务码,客户端应据此分支
	Msg       string       `json:"msg"`              // Msg 是可展示的提示,不应作为客户端的判断条件
	Data      T            `json:"data"`             // Data 返回值
	Errors    []FieldError `json:"errors,omitempty"` // Errors 仅在存在字段校验错误时返回
}

type msgOptions struct {
	msg string
}

type Option func(*msgOptions)

func WithMsg(msg string) Option {
	return func(options *msgOptions) {
		options.msg = msg
	}
}
