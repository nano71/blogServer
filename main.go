package main

import (
	"blogServer/api"
	"blogServer/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	ginServer := gin.Default()
	ginServer.Use(cors.Default())
	ginServer.Use(api.TimeoutMiddleware())
	ginServer.Static("/uploads", "./uploads")
	ginServer.Use(router.Default())

	_ = ginServer.Run()
}
