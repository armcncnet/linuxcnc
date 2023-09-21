/**
 ******************************************************************************
 * @file    config.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package Config

import (
	"armcnc/framework/utils/file"
	"github.com/gookit/color"
	"gopkg.in/ini.v1"
	"log"
	"os/exec"
)

var Get = &Data{}

type Data struct {
	Basic         DataBasic         `ini:"basic"`
	Authorization DataAuthorization `ini:"authorization"`
	Machine       DataMachine       `ini:"machine"`
}

type DataBasic struct {
	Name      string `ini:"name"`
	Version   string `ini:"version"`
	Workspace string `ini:"workspace"`
}

type DataAuthorization struct {
	Getaway string `ini:"getaway"`
	AppId   string `ini:"app_id"`
	AppKey  string `ini:"app_key"`
}

type DataMachine struct {
	Path string `ini:"path"`
}

func Init() {

	Get.Basic.Name = "armcnc"
	Get.Basic.Version = "1.0.0"
	Get.Basic.Workspace = "/opt/armcnc"

	Get.Authorization.Getaway = "https://getaway.geekros.com"

	exists, _ := FileUtils.PathExists(Get.Basic.Workspace + "/armcnc.ini")
	if !exists {
		iniFile := ini.Empty()
		err := ini.ReflectFrom(iniFile, Get)
		if err != nil {
			log.Println("[config]：" + color.Error.Sprintf("System configuration save failed"))
			return
		}

		err = iniFile.SaveTo(Get.Basic.Workspace + "/armcnc.ini")
		if err != nil {
			log.Println("[config]：" + color.Error.Sprintf("System configuration save failed"))
			return
		}
	}

	if exists {
		iniFile, err := ini.Load(Get.Basic.Workspace + "/armcnc.ini")
		if err != nil {
			log.Println("[config]：" + color.Error.Sprintf("System configuration information read failed."))
			return
		}

		var iniConfig Data
		err = iniFile.MapTo(&iniConfig)
		if err != nil {
			log.Println("[config]：" + color.Error.Sprintf("System configuration information read failed."))
			return
		}

		Get.Authorization = iniConfig.Authorization
		Get.Machine = iniConfig.Machine

		iniFile.Section("basic").Key("version").SetValue(Get.Basic.Version)
		err = iniFile.SaveTo(Get.Basic.Workspace + "/armcnc.ini")
		if err != nil {
			log.Println("[config]：" + color.Error.Sprintf("System configuration save failed"))
			return
		}

		SetEnvironment(Get.Machine.Path)
	}
}

func SetEnvironment(path string) {
	write := FileUtils.WriteFile("MACHINE_PATH="+path, Get.Basic.Workspace+"/.armcnc/environment")
	if write == nil {
		if path != "" {
			cmd := exec.Command("systemctl", "restart", "programlaunch.service")
			cmd.Output()
		} else {
			cmd := exec.Command("systemctl", "stop", "programlaunch.service")
			cmd.Output()
		}
	} else {
		cmd := exec.Command("systemctl", "stop", "programlaunch.service")
		cmd.Output()
	}
}
