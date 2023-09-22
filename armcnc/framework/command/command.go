/**
 ******************************************************************************
 * @file    command.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package Command

import (
	"armcnc/framework/command/service"
	"armcnc/framework/command/version"
	"armcnc/framework/package/config"
	"github.com/spf13/cobra"
	"os"
)

func Init() {

	command := &cobra.Command{
		Use:   "armcnc",
		Short: "Welcome to " + ConfigPackage.Get.Basic.Name + "" + ConfigPackage.Get.Basic.Version,
		Long:  "Development Team: ARMCNC https://www.armcnc.net",
	}

	command.AddCommand(VersionCommand.Start(ConfigPackage.Get.Basic.Name, ConfigPackage.Get.Basic.Version))

	command.AddCommand(ServiceCommand.Start())

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
