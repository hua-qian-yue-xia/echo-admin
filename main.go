package main

import (
	"echo-admin/pkg/rep"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()
	e.HTTPErrorHandler = rep.HTTPErrorHandler

	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/", func(c *echo.Context) error {
		return rep.OK(c, "成功")
	})

	if err := e.Start(":7001"); err != nil {
		e.Logger.Error("服务器启动失败", "error", err)
	}
}
