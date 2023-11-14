/**
 ******************************************************************************
 * @file    machine.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package MachineService

import (
	"armcnc/framework/config"
	"armcnc/framework/package/launch"
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

type responseGet struct {
	Path string              `json:"path"`
	User MachinePackage.USER `json:"user"`
	Ini  MachinePackage.INI  `json:"ini"`
}

func Get(c *gin.Context) {
	returnData := responseGet{}

	path := c.DefaultQuery("path", "")
	if path == "" {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	machine := MachinePackage.Init()
	returnData.Path = path
	returnData.User = machine.GetUser(path)
	returnData.Ini = machine.GetIni(path)

	Utils.Success(c, 0, "", returnData)
	return
}

func Set(c *gin.Context) {

	path := c.DefaultQuery("path", "")
	if path == "" {
		Utils.Error(c, 10000, "", Utils.EmptyData{})
		return
	}

	if path != Config.Get.Machine.Path {
		machine := MachinePackage.Init()
		check := machine.GetIni(path)
		if check.Emc.Version == "" {
			Utils.Error(c, 10000, "", Utils.EmptyData{})
			return
		}
		Config.Get.Machine.Path = path
		save := Config.Save()
		if !save {
			Utils.Error(c, 10000, "", Utils.EmptyData{})
			return
		}
		launch := LaunchPackage.Init()
		launch.Start(Config.Get.Machine.Path)
	}

	Utils.Success(c, 0, "", Utils.EmptyData{})
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
