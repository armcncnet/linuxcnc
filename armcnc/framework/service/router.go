/**
 ******************************************************************************
 * @file    router.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package Service

import (
	"armcnc/framework/service/config"
	"armcnc/framework/service/message"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Router() http.Handler {

	router := gin.Default()

	gin.SetMode(gin.DebugMode)

	router.Use(gin.Recovery())

	router.Use(cors.Default())

	router.Static("/files", "/opt/armcnc/files/")

	router.Static("/uploads", "/opt/armcnc/uploads/")

	config := router.Group("config")
	{
		config.GET("/index", ConfigService.Index)
	}

	message := router.Group("message")
	{
		message.GET("/service", MessageService.Service)
	}

	return router
}
