/**
 ******************************************************************************
 * @file    router.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package Service

import (
	"armcnc/framework/service/config"
	"armcnc/framework/service/machine"
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

	message := router.Group("message")
	{
		message.GET("/service", MessageService.Service)
	}

	config := router.Group("config")
	{
		config.GET("/index", ConfigService.Index)
	}

	machine := router.Group("machine")
	{
		machine.GET("/select", MachineService.Select)

		machine.GET("/get/ini/content", MachineService.GetIniContent)
	}

	return router
}
