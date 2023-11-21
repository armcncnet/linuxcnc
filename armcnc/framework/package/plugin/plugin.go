/**
 ******************************************************************************
 * @file    plugin.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package PluginPackage

import "armcnc/framework/config"

type Plugin struct {
	Path string `json:"path"`
}

type Data struct {
	Name     string `json:"name"`
	Describe string `json:"describe"`
	Script   string `json:"script"`
}

func Init() *Plugin {
	return &Plugin{
		Path: Config.Get.Basic.Workspace + "/plugins/",
	}
}

func (plugin *Plugin) Select() []Data {
	data := make([]Data, 0)
	return data
}
