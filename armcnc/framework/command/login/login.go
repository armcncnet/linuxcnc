/**
 ******************************************************************************
 * @file    login.go
 * @author  ARMCNC site:www.armcnc.net github:armcnc.github.io
 ******************************************************************************
 */

package LoginCommand

import (
	"armcnc/framework/utils"
	"armcnc/framework/utils/request"
	"bufio"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func Start() *cobra.Command {
	command := &cobra.Command{
		Use:     "login",
		Short:   "Sign in to account",
		Long:    "Sign in to account",
		Example: "armcnc login [email address]",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Println("[login]：" + color.Yellow.Text("Please enter your email address"))
				return
			}
			check := Utils.EmailValid(args[0])
			if !check {
				log.Println("[login]：" + color.Red.Text("Incorrect email forma"))
				return
			}
			request, response, _ := RequestUtils.Service("/account/login/mail/code", "GET", map[string]string{"mail": args[0]}, nil)
			if request.StatusCode != 200 {
				log.Println("[login]：" + color.Red.Text("Service request failed, please try again"))
				return
			}
			if response.Code != 0 {
				log.Println("[login]：" + color.Red.Text("Service request failed, please try again"))
				return
			}

			log.Print("[login]：" + color.Green.Text("Please enter the email verification code:"))

			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				code := scanner.Text()
				log.Println("[login]：" + color.Gray.Text(code))
			}
		},
	}
	return command
}
