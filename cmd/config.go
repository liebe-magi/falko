/*
Copyright © 2020 reeve0930 <reeve0930@gmail.com>

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
	"os"

	"github.com/spf13/cobra"
)

var (
	host string
	path string
	dest string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Change the settings",
	Run: func(cmd *cobra.Command, args []string) {
		if host == "" && path == "" && dest == "" {
			fmt.Println(conf)
			os.Exit(0)
		}
		if host != "" {
			conf.fHost = host
		}
		if path != "" {
			conf.fPath = path
		}
		if dest != "" {
			conf.cDest = dest
		}

		f, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		fmt.Println(conf)
		fmt.Fprintln(f, conf)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	configCmd.Flags().StringVarP(&host, "foltia-ip", "i", "", "set the IP address of foltia")
	configCmd.Flags().StringVarP(&path, "foltia-path", "p", "", "set the path mounted foltia")
	configCmd.Flags().StringVarP(&dest, "dest-copy", "d", "", "set the path you want to copy")
}
