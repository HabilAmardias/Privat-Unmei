package routers

import (
	"net/http"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/handlers"
	"privat-unmei/internal/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App            *gin.Engine
	StudentHandler *handlers.StudentHandlerImpl
}

func (c *RouteConfig) Setup() {
	config := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:    []string{"Content-Type", "Authorization"},
	}
	c.App.Use(cors.New(config))
	c.App.Use(middlewares.ErrorMiddleware())

	c.SetupPublicRoute()
}

func (c *RouteConfig) SetupPublicRoute() {
	v1 := c.App.Group("/api/v1")
	v1.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, dtos.Response{
			Success: true,
			Data:    "Welcome to Privat Unmei API",
		})
	})
	v1.POST("/register", c.StudentHandler.Register)
	v1.POST("/login", c.StudentHandler.Login)
}
