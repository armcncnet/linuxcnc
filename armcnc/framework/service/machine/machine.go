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

type requestUpdate struct {
	Path     string `json:"path"`
	Control  int    `json:"control"`
	Name     string `json:"name"`
	Describe string `json:"describe"`
}

func Update(c *gin.Context) {

	requestJson := requestUpdate{}
	requestData, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(requestData, &requestJson)
	if err != nil {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	machine := MachinePackage.Init()
	updateData := MachinePackage.USER{}
	updateData.Base.Name = requestJson.Name
	updateData.Base.Describe = requestJson.Describe
	updateData.Base.Control = requestJson.Control
	update := machine.UpdateUser(requestJson.Path, updateData)
	if !update {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	Utils.Success(c, 0, "", Utils.EmptyData{})
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
	returnData.Content = machine.GetIniContent(path)

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
	update := machine.UpdateIniContent(requestJson.Path, requestJson.Content)
	if !update {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	iniData := machine.GetIni(requestJson.Path)
	returnData.Machine.Path = requestJson.Path
	returnData.Machine.Version = iniData.Emc.Version
	returnData.Machine.Coordinate = iniData.Traj.Coordinates

	userData := machine.GetUser(requestJson.Path)
	returnData.Machine.Name = userData.Base.Name
	returnData.Machine.Describe = userData.Base.Describe
	returnData.Machine.Control = userData.Base.Control

	Utils.Success(c, 0, "", returnData)
	return
}
