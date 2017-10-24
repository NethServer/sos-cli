/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - nethserver@nethesis.it
 *
 * This script is part of NethServer.
 *
 * NethServer is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * NethServer is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with NethServer.  If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package cmd

import (
	"fmt"
	"errors"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"os/exec"

	"github.com/spf13/cobra"

	"sos-cli/model"
	"sos-cli/helper"
	"sos-cli/config"
)

var (
	allCloseFlag = false
)

func closeConnections() {
	resp, err := http.Get(config.API + "sessions")

	if err != nil {
		helper.RedPanic(err.Error())
	}
	defer resp.Body.Close()

	if (resp.StatusCode < 300) {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			helper.RedPanic(err.Error())
		}

		var sessions[] model.Session
		err = json.Unmarshal(body, &sessions)
		if err != nil {
			helper.RedPanic(err.Error())
		}

		for i := 0; i < len(sessions); i++ {
			closeConnection(sessions[i].SessionId)
		}
	} else {
		helper.ErrorLog("No sessions found\n")
	}
}

func closeConnection(sessionId string) {
	vpnIp := helper.GetSessionIp(sessionId)
	port := config.DEFAULT_SSH_PORT

	if (len(vpnIp) > 0) {
		helper.StartLoader()
		fmt.Printf("Try to close %s session...\n", helper.GreenString(sessionId))

		vpnCmd := exec.Command("/opt/nethsos/helpers/nethsos-stop-ssh", vpnIp, port)

		if err := vpnCmd.Start(); err != nil {
			helper.RedPanic(err.Error())
		}

		if err := vpnCmd.Wait(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				helper.RedPanic(exiterr.Error())
			} else {
				helper.RedPanic(err.Error())
			}
		}

		helper.StopLoader()
		fmt.Printf("Session %s closed!\n", helper.GreenString(sessionId))
	} else {
		helper.ErrorLog("Error: session %s not found\n", sessionId)
	}
}

var closeCmd = &cobra.Command{
	Use: "close [session-id]",
	Short: "Close Session ID and remove VPN connection",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 && !allCloseFlag {
			return errors.New(helper.RedString("requires session-id"))
		}
		return nil;
	},
	Run: func(cmd *cobra.Command, args []string) {
		if allCloseFlag {
			closeConnections()
		} else {
			closeConnection(args[0])
		}
	},
}

func init() {
	RootCmd.AddCommand(closeCmd)

	closeCmd.Flags().BoolVarP(&allCloseFlag, "all", "a", false, "Print all sessions")
}
