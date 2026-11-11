package rep

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

var ErrCommitted = errors.New("rep:HTTP响应已提交")

// committed 判断当前 HTTP 响应是否已经开始写出
func committed(c *echo.Context) bool {
	rep, _ := echo.UnwrapResponse(c.Response())
	return rep != nil && rep.Committed
}

// requestID 获取响应头中由 RequestID 中间件设置的请求 ID
func requestID(c *echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// writeOk 统一写入成功的响应
func writeOk[T any](c *echo.Context, status int, msg string, data T, options ...Option) error {
	if committed(c) {
		return ErrCommitted
	}
	settings := msgOptions{msg: msg}
	for _, option := range options {
		if option != nil {
			option(&settings)
		}
	}
	return c.JSON(status, Rep[T]{
		Code:      "OK",
		Msg:       settings.msg,
		Data:      data,
		RequestID: requestID(c),
	})
}

// OK 返回 HTTP 200和统一的成功响应
func OK[T any](c *echo.Context, data T, options ...Option) error {
	return writeOk(c, http.StatusOK, "操作成功", data, options...)
}

// Accepted 返回 HTTP 202,适用于已经受理但尚未完成的异步任务
func Accepted[T any](c *echo.Context, data T, options ...Option) error {
	return writeOk(c, http.StatusAccepted, "请求已受理", data, options...)
}

// normalizeError 将业务错误、框架错误和未知错误统一为 Error。
func normalizeError(err error) *Error {
	var statusCoder echo.HTTPStatusCoder
	if !errors.As(err, &statusCoder) {
		return New(Internal, "", WithCause(err))
	}

	// 先处理带类型的空指针，再继续解析错误链。
	switch typedErr := statusCoder.(type) {
	case *Error:
		if typedErr == nil {
			return New(Internal, "")
		}
	case *echo.HTTPError:
		if typedErr == nil {
			return New(Internal, "")
		}
	}

	// 默认 Recover 中间件捕获的 panic 统一作为内部错误处理。
	var panicErr *middleware.PanicStackError
	if errors.As(err, &panicErr) {
		return New(Internal, "", WithCause(err))
	}
	if appErr, ok := statusCoder.(*Error); ok {
		return appErr
	}

	status := statusCoder.StatusCode()
	if status < 400 || status > 599 {
		return New(Internal, "", WithCause(err))
	}
	for kind, definition := range errorDefaults {
		if definition.status == status {
			return New(Kind(kind), "", WithCause(err))
		}
	}

	msg := http.StatusText(status)
	if msg == "" {
		msg = "请求处理失败"
	}
	failure := New(Internal, msg, WithCode("HTTP_"+strconv.Itoa(status)), WithCause(err))
	failure.status = status
	return failure
}

// writeFail 统一写入失败响应。
func writeFail(c *echo.Context, err *Error) error {
	if committed(c) {
		return ErrCommitted
	}
	if c.Request().Method == http.MethodHead {
		return c.NoContent(err.StatusCode())
	}
	return c.JSON(err.StatusCode(), Rep[any]{
		RequestID: requestID(c),
		Code:      err.Code(),
		Msg:       err.Msg(),
		Data:      nil,
		Errors:    err.Fields(),
	})
}

// HTTPErrorHandler 接收 Handler 和中间件返回的错误。
// 原始请求错误由 RequestLogger 记录，这里记录响应写入失败。
func HTTPErrorHandler(c *echo.Context, err error) {
	if err == nil || committed(c) {
		return
	}
	if writeErr := writeFail(c, normalizeError(err)); writeErr != nil && !errors.Is(writeErr, ErrCommitted) {
		c.Logger().Error("写入失败响应失败",
			"error", writeErr,
			"request_id", requestID(c),
		)
	}
}
