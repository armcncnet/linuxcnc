/**
 ******************************************************************************
 * @file    machine.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package MachineCommand

import (
	"armcnc/framework/config"
	"armcnc/framework/package/launch"
	"armcnc/framework/package/machine"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"log"
)

func Start() *cobra.Command {
	command := &cobra.Command{
		Use:     "machine",
		Short:   "Machine Tool Configuration Management",
		Long:    "Machine Tool Configuration Management",
		Example: "armcnc machine [set|get] [machine_name]",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Println("[machine]：" + color.Yellow.Text("Please select the relevant operation"))
				return
			}
			if len(args) == 1 {
				if args[0] == "get" {
					if Config.Get.Machine.Path == "" {
						log.Println("[machine]: " + color.Yellow.Text("The current machine tool configuration is empty"))
						return
					}
					log.Println("[machine]: " + color.Blue.Text("The current machine tool configuration in use is: "+Config.Get.Machine.Path))
					return
				} else {
					log.Println("[machine]：" + color.Yellow.Text("Please select the relevant operation"))
					return
				}
			}
			if len(args) == 2 {
				if args[0] == "set" {
					log.Println("[machine]：" + color.Gray.Text("Please wait..."))
					if args[1] == "" {
						log.Println("[machine]：" + color.Yellow.Text("Please select the relevant operation"))
						return
					}
					machine := MachinePackage.Init()
					check := machine.GetIni(args[1])
					if check.Emc.Version == "" {
						log.Println("[machine]：" + color.Red.Text("Machine tool configuration failed. Please check and try again"))
						return
					}
					Config.Get.Machine.Path = args[1]
					save := Config.Update()
					if !save {
						log.Println("[machine]：" + color.Red.Text("Machine tool configuration failed. Please check and try again"))
						return
					}
					log.Println("[machine]: " + color.Blue.Text("The current machine tool configuration version: "+Config.Get.Machine.Path+" "+check.Emc.Version))
					launch := LaunchPackage.Init()
					launch.Start(Config.Get.Machine.Path)
				} else {
					log.Println("[machine]：" + color.Yellow.Text("Please select the relevant operation"))
					return
				}
			}
		},
	}
	return command
}
