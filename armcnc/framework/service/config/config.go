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
	Name        string `json:"name"`
	Version     string `json:"version"`
	WlanIP      string `json:"wlan_ip"`
	EthIP       string `json:"eth_ip"`
	MachinePath string `json:"machine_path"`
}

func Index(c *gin.Context) {

	returnData := responseIndex{}
	returnData.Name = Config.Get.Basic.Name
	returnData.Version = Config.Get.Basic.Version
	returnData.WlanIP = Utils.GetIP("wlan0")
	returnData.EthIP = Utils.GetIP("eth0")
	returnData.MachinePath = Config.Get.Machine.Path

	Utils.Success(c, 0, "", returnData)
	return
}
