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
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/liebe-magi/falko/db"
	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
)

// slackCmd represents the slack command
var slackCmd = &cobra.Command{
	Use:   "slack",
	Short: "Slack botを起動",
	Run: func(cmd *cobra.Command, args []string) {
		runBot()
	},
}

func init() {
	rootCmd.AddCommand(slackCmd)
}

func runBot() {
	log.Println("Slackクライアントを起動")

	api := slack.New(
		conf.sToken,
		//slack.OptionDebug(true),
		//slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	time.Sleep(1 * time.Second)

	if conf.sUser == "" || conf.sChannel == "" {
		activateSlack(rtm)
	}

	go notifyTask(rtm)

	log.Println("Slackクライアントスタンバイ完了")
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if ev.User == conf.sUser && ev.Channel == conf.sChannel {
				log.Printf("【受信】%s", ev.Msg.Text)
				text := processMsg(ev.Msg.Text)
				sendMsg(rtm, text)
			}
		}
	}
}

func sendMsg(rtm *slack.RTM, text string) {
	log.Printf("【送信】%s", strings.ReplaceAll(text, "\n", ""))
	rtm.SendMessage(rtm.NewOutgoingMessage(text, conf.sChannel))
}

func activateSlack(rtm *slack.RTM) error {
	log.Println("Slackクライアントの初期設定を開始")
	rand.Seed(time.Now().UnixNano())
	r := 0
	for {
		r = rand.Intn(10000)
		if r >= 1000 {
			break
		}
	}
	log.Printf("Slackよりこのコードを入力して下さい:%d", r)
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			num, err := strconv.Atoi(ev.Msg.Text)
			if err != nil {
				continue
			}
			if num == r {
				log.Printf("認証完了 User:%s Channel:%s", ev.Channel, ev.User)
				conf.sUser = ev.User
				conf.sChannel = ev.Channel
				writeConfig()
				return nil
			}
		}
	}
	return nil
}

func notifyTask(rtm *slack.RTM) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	err = db.InitNewAnimeDB()
	if err != nil {
		log.Fatalln(err)
	}
	h, m, err := getSlackTime()
	if err != nil {
		log.Fatalln(err)
	}
	for {
		t := time.Now()
		s := time.Date(t.Year(), t.Month(), t.Day(), h, m, 0, 0, loc)
		e := time.Date(t.Year(), t.Month(), t.Day(), h, m+1, 0, 0, loc)
		if t.After(s) && t.Before(e) {
			//予約の通知
			err = notifyReservation(rtm)
			if err != nil {
				log.Fatalln(err)
			}
			//新アニメの通知
			err = notifyNewAnime(rtm)
			if err != nil {
				log.Fatalln(err)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func notifyReservation(rtm *slack.RTM) error {
	log.Println("録画予定の通知開始")
	rl, err := getReservationList()
	if err != nil {
		return err
	}
	rl = filterReservation(rl, 1)
	text := "【本日の予約】\n"
	for _, r := range rl {
		text += fmt.Sprintf("%s [%s]\n", r.Time.Format("15:04"), r.Station)
		text += fmt.Sprintf("    %s(%d)\n", r.Title, r.EpNum)
	}
	sendMsg(rtm, text)
	return nil
}

func notifyNewAnime(rtm *slack.RTM) error {
	log.Println("新アニメ情報の通知開始")
	newAnime, err := getNewAnime()
	if err != nil {
		return fmt.Errorf("新アニメ情報の取得に失敗")
	}
	data, err := db.GetAllNewAnime()
	for _, n := range newAnime {
		exists := false
		for _, d := range data {
			if compareNewAnime(n, d) {
				exists = true
			}
		}
		if !exists {
			log.Printf("新アニメDBに追加 : %s(%d) %s", n.Title, n.TID, n.Station)
			sendMsg(rtm, makeNewAnimeInfo(n))
			time.Sleep(500 * time.Millisecond)
			err = db.InsertNewAnime(n.TID, n.Title, n.Station, n.Time)
			if err != nil {
				return err
			}
		}
	}
	data, err = db.GetAllNewAnime()
	for _, d := range data {
		exists := false
		for _, n := range newAnime {
			if compareNewAnime(n, d) {
				exists = true
			}
		}
		if !exists {
			log.Printf("新アニメDBから削除 : %s(%d) %s", d.Title, d.TID, d.Station)
			err = db.DeleteNewAnime(d.ID)
			if err != nil {
				return err
			}
		}
	}
	log.Println("新アニメ情報の通知終了")
	return nil
}

func compareNewAnime(n newAnimeInfo, d db.NewAnime) bool {
	if n.TID == d.TID && n.Station == d.Station && n.Time.Equal(d.Time) {
		return true
	}
	return false
}

func makeNewAnimeInfo(n newAnimeInfo) string {
	text := "【新アニメ情報】\n"
	text += fmt.Sprintf("%s (%d)\n", n.Title, n.TID)
	text += n.Station + " : " + n.Time.Format("2006/1/2") + " (" + jweek[n.Time.Weekday()] + ")\n"
	text += fmt.Sprintf("http://cal.syoboi.jp/tid/%d/", n.TID)
	return text
}

func processMsg(t string) string {
	tt := strings.Split(t, " ")
	if len(tt) != 2 {
		return "コマンド形式が不正"
	}
	if tt[0] != "rec" {
		return fmt.Sprintf("未定義のコマンド : %s", tt[0])
	}
	tid, err := strconv.Atoi(tt[1])
	if err != nil {
		return "TIDの指定が不正"
	}
	err = reserve(tid, 0, conf.encQuality, conf.mp2cut, conf.mp4cut)
	if err != nil {
		return "録画予約失敗"
	}
	title, _ := getTitle(tid)
	text := "【録画予約成功】\n"
	text += fmt.Sprintf("%s (%d)", title, tid)
	return text
}
