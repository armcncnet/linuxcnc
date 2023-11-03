/**
 ******************************************************************************
 * @file    launch.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package LaunchCommand

import (
	"armcnc/framework/package/launch"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"log"
)

func Start() *cobra.Command {
	command := &cobra.Command{
		Use:     "launch",
		Short:   "Service management for launch",
		Long:    "Service management for launch",
		Example: "armcnc launch [start|restart|stop]",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Println("[launch]：" + color.Yellow.Text("Please select an operation"))
				return
			}

			launch := LaunchPackage.Init()
			if args[0] == "start" {
				launch.OnStart()
			}
			if args[0] == "restart" {
				launch.OnRestart()
			}
			if args[0] == "stop" {
				launch.OnStop()
			}
			log.Println("[launch]：" + color.Green.Text("Operation successful"))
			return
		},
	}
	return command
}
