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
	"armcnc/framework/service/program"
	"armcnc/framework/service/upload"
	"armcnc/framework/service/version"
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

		machine.GET("/get", MachineService.Get)

		machine.POST("/update", MachineService.Update)

		machine.GET("/download", MachineService.Download)

		machine.GET("/delete", MachineService.Delete)

		machine.POST("/update/user", MachineService.UpdateUser)

		machine.POST("/update/launch", MachineService.UpdateLaunch)

		machine.POST("/update/hal", MachineService.UpdateHal)

		machine.POST("/update/xml", MachineService.UpdateXml)

		machine.GET("/set/current/machine", MachineService.SetCurrentMachine)

		machine.POST("/upload", UploadService.UploadMachine)
	}

	code := router.Group("program")
	{
		code.GET("/select", ProgramService.Select)

		code.GET("/read/line", ProgramService.ReadLine)

		code.GET("/read/content", ProgramService.ReadContent)

		code.POST("/update/content", ProgramService.UpdateContent)
	}

	version := router.Group("version")
	{
		version.GET("/check", VersionService.Check)
	}

	return router
}
