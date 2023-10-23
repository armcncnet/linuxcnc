/**
 ******************************************************************************
 * @file    machine.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package MachineService

import (
	"armcnc/framework/package/machine"
	"armcnc/framework/utils"
	"github.com/gin-gonic/gin"
)

type responseSelect struct {
	Machine []MachinePackage.Data `json:"machine"`
}

func Select(c *gin.Context) {

	returnData := responseSelect{}
	returnData.Machine = make([]MachinePackage.Data, 0)

	machine := MachinePackage.Init()

	returnData.Machine = machine.Select()

	Utils.Success(c, 0, "", returnData)
	return
}

type responseGetIniContent struct {
	Content string `json:"content"`
}

func GetIniContent(c *gin.Context) {

	returnData := responseGetIniContent{}

	path := c.DefaultQuery("path", "")
	if path == "" {
		Utils.Error(c, 10000, "", returnData)
		return
	}

	machine := MachinePackage.Init()
	returnData.Content = machine.GetContent(path)

	Utils.Success(c, 0, "", returnData)
	return
}
