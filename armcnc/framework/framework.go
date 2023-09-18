/**
 ******************************************************************************
 * @file    framework.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package Framework

import (
	"armcnc/framework/command"
	"armcnc/framework/config"
)

func Init() {

	// 初始化全局配置
	Config.Init()

	// 初始化命令行工具
	Command.Init()
}
