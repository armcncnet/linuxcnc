/**
 ******************************************************************************
 * @file    backup.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package BackupPackage

import "armcnc/framework/config"

type Backup struct {
	Path string `json:"path"`
}

type Data struct {
}

func Init() *Backup {
	return &Backup{
		Path: Config.Get.Basic.Workspace + "/backups/",
	}
}
