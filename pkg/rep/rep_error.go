package rep

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
)

type Kind uint8

const (
	// Internal 表示未预期的服务端错误
	Internal Kind = iota
	// InvalidArgument 表示请求参数无效
	InvalidArgument
	// Unauthenticated 表示调用方未通过身份认证
	Unauthenticated
	// PermissionDenied 表示调用方没有操作权限
	PermissionDenied
	// NotFound 表示目标资源不存在
	NotFound
	// Conflict 表示资源状态或唯一性冲突
	Conflict
	// ValidationFailed 表示业务字段校验失败
	ValidationFailed
	// RateLimited 表示请求超过允许频率
	RateLimited
	// Unavailable 表示服务暂时不可用
	Unavailable
	// DeadlineExceeded 表示等待上游服务响应超时
	DeadlineExceeded
)

func defaults(kind Kind) Error {
	if int(kind) >= len(errorDefaults) {
		kind = Internal
	}
	err := errorDefaults[kind]
	err.kind = kind
	return err
}

func (k Kind) DefaultCode() string {
	return defaults(k).code
}

func (k Kind) DefaultMsg() string {
	return defaults(k).msg
}

type FieldError struct {
	Field string `json:"field"` // Field 字段名
	Code  string `json:"code"`  // Code 可选的字段校验规则码
	Msg   string `json:"msg"`   // Msg 向客户端展示的校验提示
}

type Error struct {
	kind   Kind         // kind 错误分类,用于映射 HTTP 状态
	status int          // status HTTP 状态码
	code   string       // code 稳定、具体的业务码，供前端判断
	msg    string       // msg 面向用户的提示
	cause  error        // cause 原始错误,供日志和错误链分析
	fields []FieldError // fields 字段级校验信息
}

// errorDefaults 集中定义每种错误的 HTTP 状态、业务码和默认文案。
// New 按值复制默认错误，再应用当前请求的选项。
var errorDefaults = [...]Error{
	Internal: {
		status: http.StatusInternalServerError,
		code:   "INTERNAL_ERROR",
		msg:    "服务器内部错误，请稍后重试",
	},
	InvalidArgument: {
		status: http.StatusBadRequest,
		code:   "INVALID_ARGUMENT",
		msg:    "请求参数不正确",
	},
	Unauthenticated: {
		status: http.StatusUnauthorized,
		code:   "UNAUTHENTICATED",
		msg:    "请先登录",
	},
	PermissionDenied: {
		status: http.StatusForbidden,
		code:   "PERMISSION_DENIED",
		msg:    "无权执行此操作",
	},
	NotFound: {
		status: http.StatusNotFound,
		code:   "NOT_FOUND",
		msg:    "资源不存在",
	},
	Conflict: {
		status: http.StatusConflict,
		code:   "CONFLICT",
		msg:    "资源状态冲突",
	},
	ValidationFailed: {
		status: http.StatusUnprocessableEntity,
		code:   "VALIDATION_FAILED",
		msg:    "参数校验失败",
	},
	RateLimited: {
		status: http.StatusTooManyRequests,
		code:   "RATE_LIMITED",
		msg:    "请求过于频繁，请稍后重试",
	},
	Unavailable: {
		status: http.StatusServiceUnavailable,
		code:   "SERVICE_UNAVAILABLE",
		msg:    "服务暂时不可用，请稍后重试",
	},
	DeadlineExceeded: {
		status: http.StatusGatewayTimeout,
		code:   "DEADLINE_EXCEEDED",
		msg:    "请求处理超时，请稍后重试",
	},
}

type ErrorOption func(*Error)

// New 创建一个独立的业务错误实例
func New(kind Kind, msg string, errOptions ...ErrorOption) *Error {
	err := defaults(kind)
	if strings.TrimSpace(msg) != "" {
		err.msg = msg
	}
	for _, option := range errOptions {
		if option != nil {
			option(&err)
		}
	}
	return &err
}

func WithCode(code string) ErrorOption {
	code = strings.TrimSpace(code)
	return func(err *Error) {
		if code != "" {
			err.code = code
		}
	}
}

func WithCause(cause error) ErrorOption {
	return func(err *Error) {
		err.cause = cause
	}
}

func WithFields(fields ...FieldError) ErrorOption {
	clone := slices.Clone(fields)
	return func(err *Error) {
		err.fields = clone
	}
}

func (err *Error) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.cause
}

func (err *Error) Kind() Kind {
	return err.kind
}

// StatusCode 实现 Echo 的 HTTPStatusCoder 接口。
func (err *Error) StatusCode() int {
	if err == nil {
		return defaults(Internal).status
	}
	if err.status == 0 {
		return defaults(err.kind).status
	}
	return err.status
}

func (err *Error) Code() string {
	if err.code == "" {
		return err.kind.DefaultCode()
	}
	return err.code
}

func (err *Error) Msg() string {
	if err.msg == "" {
		return err.kind.DefaultMsg()
	}
	return err.msg
}

func (err *Error) Fields() []FieldError {
	return slices.Clone(err.fields)
}

// Error 返回用于日志记录的错误信息。
func (err *Error) Error() string {
	if err == nil {
		return "<nil>"
	}
	if err.cause != nil {
		return fmt.Sprintf(
			"code=%s,msg=%s,cause=%v",
			err.Code(),
			err.Msg(),
			err.cause,
		)
	}
	return fmt.Sprintf(
		"code=%s,message=%s",
		err.Code(),
		err.Msg(),
	)
}
