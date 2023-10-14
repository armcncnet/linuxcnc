/**
 ******************************************************************************
 * @file    launch.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package LaunchPackage

import (
	"armcnc/framework/config"
	"armcnc/framework/utils/file"
	"os/exec"
)

type Launch struct {
	Path string `json:"path"`
}

type Data struct {
}

func Init() *Launch {
	return &Launch{
		Path: Config.Get.Basic.Workspace,
	}
}

func (launch *Launch) Start(machine string) {
	write := FileUtils.WriteFile("MACHINE_PATH="+machine, Config.Get.Basic.Workspace+"/.armcnc/environment")
	if write == nil {
		if machine != "" {
			cmd := exec.Command("systemctl", "restart", "armcnc_launch.service")
			cmd.Output()
		} else {
			cmd := exec.Command("systemctl", "stop", "armcnc_launch.service")
			cmd.Output()
		}
	} else {
		cmd := exec.Command("systemctl", "stop", "armcnc_launch.service")
		cmd.Output()
	}
}

func (launch *Launch) Stop() {
	cmd := exec.Command("systemctl", "stop", "armcnc_launch.service")
	cmd.Output()
}
