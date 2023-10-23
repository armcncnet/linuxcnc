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
	"github.com/goccy/go-json"
	"io/ioutil"
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

type requestUpdateIniContent struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type responseUpdateIniContent struct {
	Machine MachinePackage.Data `json:"machine"`
}

func UpdateIniContent(c *gin.Context) {

	returnData := responseUpdateIniContent{}

	requestJson := requestUpdateIniContent{}

	requestData, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(requestData, &requestJson)
	if err != nil {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	machine := MachinePackage.Init()
	update := machine.UpdateContent(requestJson.Path, requestJson.Content)
	if !update {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	machineData := machine.Get(requestJson.Path)
	returnData.Machine.Path = requestJson.Path
	returnData.Machine.Name = machineData.EMC.MACHINE
	returnData.Machine.Describe = machineData.EMC.DESCRIBE
	returnData.Machine.Version = machineData.EMC.VERSION
	returnData.Machine.ControlType = machineData.EMC.CONTROL_TYPE
	returnData.Machine.Coordinate = machineData.TRAJ.COORDINATES

	Utils.Success(c, 0, "", returnData)
	return
}
