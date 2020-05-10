/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/reeve0930/foltia/db"
	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset copy status",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalln("Please set TID:EpNum")
		} else if len(args) == 1 {
			t, e, err := parseTIDEpNum(args[0])
			if err != nil {
				log.Fatalln(err)
			}
			err = resetCopyStatus(t, e)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Fatalln("Args are too many")
		}
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
	// resetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func parseTIDEpNum(s string) (int, int, error) {
	ss := strings.Split(s, ":")
	if len(ss) != 2 {
		return 0, 0, fmt.Errorf("Please check args format -> TID:EpNum")
	}
	t, err := strconv.Atoi(ss[0])
	if err != nil {
		return 0, 0, err
	}
	e, err := strconv.Atoi(ss[1])
	if err != nil {
		return 0, 0, err
	}
	return t, e, nil
}

func resetCopyStatus(t int, e int) error {
	data, err := db.GetAllEpisode()
	if err != nil {
		return err
	}
	for _, d := range data {
		if d.TID == t && d.EpNum == e {
			title, err := getTitle(t)
			if err != nil {
				return err
			}
			log.Printf("Reset copy status : (%d)%s (%d:%s)", d.TID, title, d.EpNum, d.EpTitle)
			err = db.UpdateEpisode(d.ID, d.TID, d.EpNum, d.EpTitle, false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
