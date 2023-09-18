/**
 ******************************************************************************
 * @file    config.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package Config

var Get = &Data{}

type Data struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Workspace string `json:"workspace"`
}

func Init() {

	Get.Name = "armcnc"
	Get.Version = "1.0.0"
	Get.Workspace = "/opt/armcnc"

}
