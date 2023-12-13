/**
 ******************************************************************************
 * @file    chromium.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package ChromiumPackage

import "armcnc/framework/config"

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
