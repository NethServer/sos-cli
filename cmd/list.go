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
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"strconv"
	"bytes"

	"github.com/spf13/cobra"

	"sos-cli/model"
	"sos-cli/config"
	"sos-cli/helper"
)

var (
	jsonFlag = false
)

func printJSON(body []byte) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		helper.RedPanic(err.Error())
	}
	fmt.Println(string(prettyJSON.Bytes()))
}

func printSession(session model.SessionOutput) {
	if (jsonFlag) {
		jsonPrint := []byte(`{
			"session":"` +session.SessionId + `",
			"lk":"` + session.Lk + `",
			"vpn":"` + session.VpnIp + `",
			"started":"` + session.Started + `"
		}`)
		printJSON(jsonPrint);
	} else {
		fmt.Printf("session: %s\n", helper.GreenString(session.SessionId))
		fmt.Printf("  lk:\t\t%s\n", session.Lk)
		fmt.Printf("  vpn:\t\t%s\n", session.VpnIp)
		fmt.Printf("  started:\t%s\n\n", session.Started)
	}
}

func listSessions() {
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
			sessionId := sessions[i].SessionId
			vpn := sessions[i].VpnIp
			lk := sessions[i].Lk

			started, err := strconv.ParseInt(strconv.Itoa(sessions[i].Started), 10, 64)
			if err != nil {
				helper.RedPanic(err.Error())
			}

			sessionToPrint := model.SessionOutput{
				SessionId: sessionId,
				Lk: lk,
				VpnIp: vpn,
				Started: time.Unix(started, 0).String(),
			}

			printSession(sessionToPrint)

		}
	} else {
		helper.ErrorLog("No sessions found\n")
	}
}

func listSession(sessionId string) {
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

		vpn := session.VpnIp
		lk := session.Lk

		started, err := strconv.ParseInt(strconv.Itoa(session.Started), 10, 64)
		if err != nil {
			helper.RedPanic(err.Error())
		}

		sessionToPrint := model.SessionOutput{
			SessionId: sessionId,
			Lk: lk,
			VpnIp: vpn,
			Started: time.Unix(started, 0).String(),
		}

		printSession(sessionToPrint)
	} else {
		helper.ErrorLog("No session %s found\n", sessionId)
	}
}

var listCmd = &cobra.Command{
	Use: "list [session-id]",
	Short: "Show all VPNs of connected servers",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			listSession(args[0])
		} else {
			listSessions()
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&jsonFlag, "json", "j", false, "Print output in JSON format")
}
