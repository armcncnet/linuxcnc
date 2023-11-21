/**
 ******************************************************************************
 * @file    backup.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package BackupPackage

import (
	"armcnc/framework/config"
	"github.com/djherbis/times"
	"os"
)

type Backup struct {
	Path string `json:"path"`
}

type Data struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Data string `json:"data"`
}

func Init() *Backup {
	return &Backup{
		Path: Config.Get.Basic.Workspace + "/backups/",
	}
}

func (backup *Backup) Select() []Data {
	data := make([]Data, 0)

	files, err := os.ReadDir(backup.Path)
	if err != nil {
		return data
	}

	for _, file := range files {
		item := Data{}
		item.Name = file.Name()
		item.Path = file.Name()
		timeData, _ := times.Stat(backup.Path + file.Name())
		createTime := timeData.BirthTime()
		item.Data = createTime.Format("2006-01-02 15:04:05")
		data = append(data, item)
	}

	return data
}
