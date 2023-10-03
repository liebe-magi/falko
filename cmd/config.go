/*
Copyright © 2020 liebe-magi <liebe.magi@gmail.com>

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
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	host         string
	path         string
	dest         string
	filename     string
	filetype     string
	dropThresh   int
	encQuality   int
	mp2cut       int
	mp4cut       int
	slackToken   string
	slackTime    string
	slackName    string
	slackChannel string
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
		if slackToken != "" {
			conf.sToken = slackToken
		}
		err := checkTime(slackTime)
		if err != nil {
			log.Fatalln(err)
		} else {
			if slackTime != "00:00" {
				conf.sTime = slackTime
			}
		}

		err = writeConfig()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(conf)
	},
}

func writeConfig() error {
	f, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, conf)
	return nil
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
	configCmd.Flags().StringVarP(&slackToken, "slack_token", "b", "", "Slack botトークンの設定")
	configCmd.Flags().StringVarP(&slackTime, "slack_time", "c", "00:00", "Slack通知を送る時間の設定")
}

func checkFlags() bool {
	if host == "" && path == "" && dest == "" && filename == "" && filetype == "" && dropThresh == 0 && encQuality == -1 && mp2cut == -1 && mp4cut == -1 && slackToken == "" && slackTime == "00:00" {
		return true
	}
	return false
}

func checkTime(t string) error {
	tt := strings.Split(t, ":")
	if len(tt) != 2 {
		return fmt.Errorf("時刻指定が不正値 : slack_time")
	}
	h, err := strconv.Atoi(tt[0])
	if err != nil {
		return nil
	}
	if h < 0 || h > 23 {
		return fmt.Errorf("時刻の範囲が不正")
	}
	m, err := strconv.Atoi(tt[1])
	if err != nil {
		return nil
	}
	if m < 0 || m > 59 {
		return fmt.Errorf("時刻の範囲が不正")
	}
	return nil
}

func getSlackTime() (int, int, error) {
	tt := strings.Split(conf.sTime, ":")
	h, err := strconv.Atoi(tt[0])
	if err != nil {
		return 0, 0, err
	}
	m, err := strconv.Atoi(tt[1])
	if err != nil {
		return 0, 0, err
	}
	return h, m, nil
}
