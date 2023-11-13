/**
 ******************************************************************************
 * @file    version.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package VersionService

import (
	"armcnc/framework/package/version"
	"armcnc/framework/utils"
	"github.com/gin-gonic/gin"
)

type responseCheck struct {
	ARMCNC   string `json:"armcnc"`
	LINUXCNC string `json:"linuxcnc"`
	SDK      string `json:"sdk"`
}

func Check(c *gin.Context) {
	returnData := responseCheck{}

	version := VersionPackage.Init()
	data := version.Get()

	returnData.ARMCNC = data.ARMCNC
	returnData.LINUXCNC = data.LINUXCNC
	returnData.SDK = data.SDK

	Utils.Success(c, 0, "", returnData)
	return
}
