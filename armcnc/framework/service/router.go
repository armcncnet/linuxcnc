/**
 ******************************************************************************
 * @file    router.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package Service

import (
	"armcnc/framework/config"
	"armcnc/framework/service/backup"
	"armcnc/framework/service/config"
	"armcnc/framework/service/machine"
	"armcnc/framework/service/message"
	"armcnc/framework/service/plugin"
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

	router.Static("/backup", Config.Get.Basic.Workspace+"/backups/")

	router.Static("/programs", Config.Get.Basic.Workspace+"/programs/")

	router.Static("/uploads", Config.Get.Basic.Workspace+"/uploads/")

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

	program := router.Group("program")
	{
		program.GET("/select", ProgramService.Select)

		program.GET("/read/line", ProgramService.ReadLine)

		program.GET("/read/content", ProgramService.ReadContent)

		program.GET("/download", ProgramService.Download)

		program.POST("/update/content", ProgramService.UpdateContent)

		program.GET("/delete", ProgramService.Delete)

		program.POST("/upload", UploadService.UploadProgram)
	}

	plugin := router.Group("plugin")
	{
		plugin.GET("/select", PluginService.Select)
	}

	settings := router.Group("settings")
	{
		settings.GET("/backup/select", BackupService.Select)

		settings.GET("/backup/pack", BackupService.Pack)

		settings.GET("/backup/delete", BackupService.Delete)

		settings.GET("/version/check", VersionService.Check)
	}

	return router
}
