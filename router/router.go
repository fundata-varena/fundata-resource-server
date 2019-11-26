package router

import (
	"github.com/fundata-varena/fundata-resource-server/controller"
	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
}

func NewRouter() *Router {
	engine := gin.Default()

	engine.GET("/resource", controller.GetResource)
	engine.GET("/resources", controller.GetResources)

	return &Router{engine}
}
