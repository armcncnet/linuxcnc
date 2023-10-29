/**
 ******************************************************************************
 * @file    service.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package ServiceCommand

import (
	"armcnc/framework/config"
	"armcnc/framework/package/launch"
	"armcnc/framework/package/machine"
	"armcnc/framework/service"
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var Get = &Data{}

type Data struct {
	Group errgroup.Group
}

func Start() *cobra.Command {
	command := &cobra.Command{
		Use:     "service",
		Short:   "Start core service",
		Long:    "Start core service",
		Example: "armcnc service",
		Hidden:  true,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("[service]: " + color.Gray.Text("Core service is starting..."))
			start := &http.Server{
				Addr:           fmt.Sprintf(":%d", Config.Get.Basic.Port),
				Handler:        Service.Router(),
				ReadTimeout:    60 * time.Second,
				WriteTimeout:   60 * time.Second,
				MaxHeaderBytes: 1 << 20,
			}

			Get.Group.Go(func() error {
				return start.ListenAndServe()
			})

			launch := LaunchPackage.Init()

			if Config.Get.Machine.Path != "" {
				machine := MachinePackage.Init()
				check := machine.GetIni(Config.Get.Machine.Path)
				if check.Emc.Version != "" {
					launch.Start(Config.Get.Machine.Path)
				}
			}

			log.Println("[service]: " + color.Info.Sprintf("Core service started successfully"))

			quit := make(chan os.Signal)
			signal.Notify(quit, os.Interrupt)
			<-quit

			launch.Stop()
			log.Println("[service]: " + color.Info.Sprintf("Core service exit"))

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := start.Shutdown(ctx); err != nil {
			}

			if err := Get.Group.Wait(); err != nil {
			}
		},
	}
	return command
}
