/**
 ******************************************************************************
 * @file    chromium.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package ChromiumPackage

import (
	"armcnc/framework/config"
	"os/exec"
)

type Chromium struct {
	Path string `json:"path"`
}

type Data struct {
}

func Init() *Chromium {
	return &Chromium{
		Path: Config.Get.Basic.Workspace + "/",
	}
}

func (chromium *Chromium) OnStart() {
	cmd := exec.Command("systemctl", "start", "armcnc_chromium.service")
	cmd.Output()
}

func (chromium *Chromium) OnReStart() {
	cmd := exec.Command("systemctl", "restart", "armcnc_chromium.service")
	cmd.Output()
}

func (chromium *Chromium) OnStop() {
	cmd := exec.Command("systemctl", "stop", "armcnc_chromium.service")
	cmd.Output()
}
