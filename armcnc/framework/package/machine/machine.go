/**
 ******************************************************************************
 * @file    machine.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package MachinePackage

import "armcnc/framework/package/config"

type Machine struct {
	Path string `json:"path"`
}

type Data struct {
}

func Init() *Machine {
	return &Machine{
		Path: ConfigPackage.Get.Basic.Workspace + "/configs/",
	}
}

func (manager *Machine) Select() []Data {
	data := make([]Data, 0)
	return data
}
