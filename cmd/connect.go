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
	"os/exec"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/spf13/cobra"

	"sos-cli/model"
	"sos-cli/config"
	"sos-cli/helper"
)

func getSessionIP(sessionId string) string {
	resp, err := http.Get(config.API + "sessions/" + sessionId)

	if err != nil {
		helper.RedPanic(err.Error())
	}
	defer resp.Body.Close()

	if (resp.StatusCode < 300) {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			helper.RedPanic(err.Error())
		}

		var session model.Session
		err = json.Unmarshal(body, &session)
		if err != nil {
			helper.RedPanic(err.Error())
		}

		return session.VpnIp
	} else {
		return ""
	}
}

func connectSession(sessionId string) {
	vpnIp := getSessionIP(sessionId)

	if (len(vpnIp) > 0) {
		helper.StartLoader()
		fmt.Printf("Try connection on %s session...\n", helper.GreenString(sessionId))

		vpmCmd := exec.Command("echo", vpnIp)

		if err := vpmCmd.Start(); err != nil {
			helper.RedPanic(err.Error())
		}

		if err := vpmCmd.Wait(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				helper.RedPanic(exiterr.Error())
			} else {
				helper.RedPanic(err.Error())
			}
		}

		helper.StopLoader()
		fmt.Printf("Connection on %s session established!\n", helper.GreenString(sessionId))
	} else {
		helper.ErrorLog("Error: session %s not found\n", sessionId)
	}
}

var connectCmd = &cobra.Command{
	Use: "connect [session-id]",
	Short: "Connect to server by specify Session ID",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
		  return errors.New("requires session-id")
		}
		return nil;
	},
	Run: func(cmd *cobra.Command, args []string) {
		sessionId := args[0]

		connectSession(sessionId)
	},
}

func init() {
	RootCmd.AddCommand(connectCmd)
}
