/**
 ******************************************************************************
 * @file    version.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package VersionPackage

import (
	"armcnc/framework/config"
	"bufio"
	"os/exec"
	"regexp"
	"strings"
)

type Version struct {
	ARMCNC   string `json:"armcnc"`
	LINUXCNC string `json:"linuxcnc"`
	SDK      string `json:"sdk"`
}

func Init() *Version {
	return &Version{
		ARMCNC:   "",
		LINUXCNC: "",
		SDK:      "",
	}
}

func (version *Version) Get() Version {
	data := Version{}
	data.ARMCNC = version.ArmCNC()
	data.LINUXCNC = version.LinuxCNC()
	data.SDK = version.Python()
	return data
}

func (version *Version) ArmCNC() string {
	data := Config.Get.Basic.Version
	return data
}

func (version *Version) Python() string {
	data := ""
	output, err := exec.Command("bash", "-c", "python3 -B -c \"import pkg_resources; print(pkg_resources.get_distribution('armcnc').version)\"").Output()
	if err == nil {
		data = strings.TrimSpace(string(output))
	}
	return data
}

func (version *Version) LinuxCNC() string {
	data := ""
	output, err := exec.Command("bash", "-c", "dpkg -l | grep linuxcnc").Output()
	if err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(output)))
		versionRegex := regexp.MustCompile(`\d+\.\d+\.\d+`)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "linuxcnc-uspace") {
				matches := versionRegex.FindStringSubmatch(line)
				if len(matches) > 0 {
					data = matches[0]
				}
			}
		}
	}
	return data
}
