/**
 ******************************************************************************
 * @file    config.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package ConfigService

import (
	"armcnc/framework/config"
	"armcnc/framework/utils"
	"github.com/gin-gonic/gin"
)

type responseIndex struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func Index(c *gin.Context) {

	returnData := responseIndex{}
	returnData.Name = Config.Get.Basic.Name
	returnData.Version = Config.Get.Basic.Version

	Utils.Success(c, 0, "", returnData)
	return
}
