/**
 ******************************************************************************
 * @file    version.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package VersionCommand

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"log"
)

func Start(name string, version string) *cobra.Command {
	command := &cobra.Command{
		Use:     "version",
		Short:   "Get version number",
		Long:    "Get version number",
		Example: "armcnc version",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("[version]: " + color.Green.Text(name+" "+version+" "))
		},
	}
	return command
}
