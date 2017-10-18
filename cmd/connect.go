// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
		fmt.Printf("Try connection on %s session...\n", sessionId)

		vpmCmd := exec.Command("sleep", vpnIp)

		if err := vpmCmd.Start(); err != nil {
			helper.RedPanic(err.Error())
		}

		if err := vpmCmd.Wait(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				helper.RedPanic(exiterr.Error())
			} else {
				fmt.Printf("vpmCmd.Wait: %v", err)
			}
		}

		helper.StopLoader()
		helper.SuccessLog("Connection on %s session established!\n", sessionId)
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
