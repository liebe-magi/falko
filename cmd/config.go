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
	host       string
	path       string
	dest       string
	filename   string
	filetype   string
	dropThresh int
	encQuality int
	mp2cut     int
	mp4cut     int
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "設定変更",
	Run: func(cmd *cobra.Command, args []string) {
		if checkFlags() {
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
		if filename != "" {
			conf.cFilename = filename
		}
		if filetype != "" {
			conf.cFiletype = filetype
		}
		if dropThresh >= 0 {
			conf.cDropThresh = dropThresh
		}
		if encQuality >= 0 {
			conf.encQuality = encQuality
		}
		if mp2cut >= 0 {
			conf.mp2cut = mp2cut
		}
		if mp4cut >= 0 {
			conf.mp4cut = mp4cut
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

	configCmd.Flags().StringVarP(&host, "foltia-ip", "i", "", "foltia ANIME LOCKERのIPアドレスを設定")
	configCmd.Flags().StringVarP(&path, "foltia-path", "s", "", "foltia ANIME LOCKERをマウントしているディレクトリを設定")
	configCmd.Flags().StringVarP(&dest, "dest-copy", "d", "", "コピー先のディレクトリを設定")
	configCmd.Flags().StringVarP(&filename, "filename", "n", "", "コピー時のファイル名フォーマットを設定")
	configCmd.Flags().StringVarP(&filetype, "file-type", "t", "", "コピーするファイルタイプを設定")
	configCmd.Flags().IntVarP(&dropThresh, "drop-thresh", "r", -1, "コピー時のTSドロップ数の閾値設定")
	configCmd.Flags().IntVarP(&encQuality, "encode-quality", "e", -1, "予約時のエンコード設定")
	configCmd.Flags().IntVarP(&mp2cut, "mp2cm_cut", "x", -1, "予約時のMPEG2編集設定")
	configCmd.Flags().IntVarP(&mp4cut, "mp4cm_cut", "y", -1, "予約時のMP4編集設定")
}

func checkFlags() bool {
	if host == "" && path == "" && dest == "" && filename == "" && filetype == "" && dropThresh == 0 {
		return true
	}
	return false
}
