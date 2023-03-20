package http

import (
	"github.com/labstack/echo/v4"
	"goWebService/handler"
)

func AccountRoutes(accountGroup *echo.Group, h handler.Handlers) {
	accountGroup.POST("/save", h.Save())
	accountGroup.PUT("/update", h.Update())
	accountGroup.DELETE("/delete/:id", h.Delete())
	accountGroup.GET("/get/:id", h.GetByID())
}
