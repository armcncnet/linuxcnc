/**
 ******************************************************************************
 * @file    command.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package Command

import (
	"armcnc/framework/config"
	"github.com/spf13/cobra"
	"os"
)

func Init() {

	command := &cobra.Command{
		Use:   "armcnc",
		Short: "Welcome to " + Config.Get.Name,
		Long:  "Development Team: ARMCNC https://www.armcnc.net",
	}

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
