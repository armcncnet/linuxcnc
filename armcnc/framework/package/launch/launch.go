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
	"os"
	"os/exec"
	"path/filepath"
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
	write := FileUtils.WriteFile("MACHINE_PATH="+machine, "/tmp/environment")
	if write == nil {
		if machine != "" {
			launch.OnRestart()
		} else {
			launch.OnStop()
		}
	} else {
		launch.OnStop()
	}
}

func (launch *Launch) OnStart() {
	exists, _ := FileUtils.PathExists("/tmp/linuxcnc.lock")
	if !exists {
		launch.OnRemoveTmp()
		cmd := exec.Command("systemctl", "start", "armcnc_linuxcnc.service")
		cmd.Output()
		cmd = exec.Command("systemctl", "start", "armcnc_launch.service")
		cmd.Output()
	}
}

func (launch *Launch) OnRestart() {
	launch.OnRemoveTmp()
	cmd := exec.Command("systemctl", "restart", "armcnc_linuxcnc.service")
	cmd.Output()
	cmd = exec.Command("systemctl", "restart", "armcnc_launch.service")
	cmd.Output()
}

func (launch *Launch) OnStop() {
	cmd := exec.Command("systemctl", "stop", "armcnc_launch.service")
	cmd.Output()
	cmd = exec.Command("systemctl", "stop", "armcnc_linuxcnc.service")
	cmd.Output()
}

func (launch *Launch) OnRemoveTmp() {
	files, err := filepath.Glob(filepath.Join("/tmp/", "linuxcnc.*"))
	if err == nil {
		for _, file := range files {
			os.Remove(file)
		}
	}
}
