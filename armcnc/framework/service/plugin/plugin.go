/**
 ******************************************************************************
 * @file    plugin.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package PluginService

import (
	"armcnc/framework/package/plugin"
	"armcnc/framework/utils"
	"github.com/gin-gonic/gin"
)

type responseSelect struct {
	Plugin []PluginPackage.Data `json:"plugin"`
}

func Select(c *gin.Context) {

	returnData := responseSelect{}
	returnData.Plugin = make([]PluginPackage.Data, 0)

	plugin := PluginPackage.Init()
	returnData.Plugin = plugin.Select()

	Utils.Success(c, 0, "", returnData)
	return
}
