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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/reeve0930/foltia/db"
	"github.com/spf13/cobra"
)

var jweek = [7]string{"日", "月", "火", "水", "木", "金", "土"}

type newAnimeInfo struct {
	TID     int
	Station string
	Title   string
	Time    time.Time
}

func (n newAnimeInfo) String() string {
	return fmt.Sprintf("%s : %s(%d) %s",
		n.Time.Format("2006/01/02")+"("+jweek[n.Time.Weekday()]+") "+n.Time.Format("15:04"),
		n.Title,
		n.TID,
		n.Station,
	)
}

type foltiaStatus struct {
	version        string
	serial         string
	storage        string
	storageRemain  string
	storagePercent int
	runningDays    int
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "foltia ANIME LOCKERの情報を確認",
	Run: func(cmd *cobra.Command, args []string) {
		newAnime, err := cmd.Flags().GetBool("new-anime")
		if err != nil {
			log.Fatalln(err)
		}
		tidFlag, err := cmd.Flags().GetBool("tid")
		if err != nil {
			log.Fatalln(err)
		}
		packet, err := cmd.Flags().GetInt("packet")
		if err != nil {
			log.Fatalln(err)
		}
		c := checkFlag(newAnime, tidFlag, packet)
		if c == 0 {
			err = showStatus()
			if err != nil {
				log.Fatalln(err)
			}
		} else if c == 1 {
			if newAnime {
				err = checkNewAnime()
				if err != nil {
					log.Fatalln(err)
				}
			}
			if tidFlag {
				err = showTitle()
				if err != nil {
					log.Fatalln(err)
				}
			}
			if packet != 0 {
				err = showTSInfo(packet)
				if err != nil {
					log.Fatalln(err)
				}
			}
		} else {
			log.Println("複数のフラグを同時に指定できません")
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().BoolP("new-anime", "n", false, "新アニメリストの確認")
	checkCmd.Flags().BoolP("tid", "t", false, "TIDの一覧を確認")
	checkCmd.Flags().IntP("packet", "p", 0, "TSドロップの発生しているファイルを確認 (閾値を指定)")
}

func checkFlag(n bool, t bool, p int) int {
	count := 0
	if n {
		count++
	}
	if t {
		count++
	}
	if p != 0 {
		count++
	}
	return count
}

func showTSInfo(p int) error {
	data, err := db.GetAllVideoFile()
	if err != nil {
		return err
	}
	title, err := db.GetAllTitle()
	if err != nil {
		return err
	}
	for _, d := range data {
		if d.Drop > p {
			var name string
			for _, t := range title {
				if t.TID == d.TID {
					name = t.Title
					break
				}
			}
			if name != "" {
				fmt.Printf("%s(%d:%d) (%d) D:%d S:%d\n", name, d.TID, d.EpNum, d.PID, d.Drop, d.Scramble)
			} else {
				return fmt.Errorf("TIDが見つかりません : %d", d.TID)
			}
		}
	}
	return nil
}

func showStatus() error {
	s, err := getStatus()
	if err != nil {
		return err
	}
	fmt.Println("foltia ANIME LOCKERシステム情報")
	fmt.Printf("  Version : %s\n", s.version)
	fmt.Printf("  Serial No. : %s\n", s.serial)
	fmt.Printf("  Running : %d days\n", s.runningDays)
	fmt.Printf("  Storage : %s/%s (Rem %d%s)\n", s.storageRemain, s.storage, 100-s.storagePercent, "%")
	return nil
}

func getStatus() (foltiaStatus, error) {
	url := "http://" + conf.fHost + "/setup/about.php"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return foltiaStatus{}, err
	}
	var fs foltiaStatus
	fs.version = strings.TrimSpace(doc.Find("#setUpTable > table > tbody > tr:nth-child(2) > td:nth-child(2)").Text())
	fs.serial = strings.TrimSpace(strings.Split(doc.Find("#setUpTable > table > tbody > tr:nth-child(3) > td:nth-child(2)").Text(), "\n")[0])
	fs.runningDays, err = strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(doc.Find("#setUpTable > table > tbody > tr:nth-child(8) > td:nth-child(2)").Text(), "days")))
	if err != nil {
		return foltiaStatus{}, err
	}
	url = "http://" + conf.fHost + "/recorded/recfiles_tid.php"
	doc, err = goquery.NewDocument(url)
	if err != nil {
		return foltiaStatus{}, err
	}
	fs.storagePercent, err = strconv.Atoi(strings.TrimSuffix(doc.Find("#HDDremainder > dl > dd > span.spent").Text(), "%"))
	if err != nil {
		return foltiaStatus{}, err
	}
	fs.storage = strings.TrimSpace(doc.Find("#HDDtotal").Text())
	fs.storageRemain = strings.TrimSpace(doc.Find("#HDDrest").Text())
	return fs, nil
}

func checkNewAnime() error {
	nail, err := getNewAnime()
	if err != nil {
		return err
	}
	for _, n := range nail {
		fmt.Println(n)
	}
	return nil
}

func getNewAnime() ([]newAnimeInfo, error) {
	url := "http://" + conf.fHost + "/animeprogram/index.php?filter=crp&view=np"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return []newAnimeInfo{}, err
	}
	var nail []newAnimeInfo
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return []newAnimeInfo{}, err
	}
	doc.Find("#contents > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		if i != 0 {
			var n newAnimeInfo
			n.TID, _ = strconv.Atoi(s.Find("td[class=TID]").Text())
			n.Station = s.Find("td[class=station]").Text()
			n.Title = s.Find("td[class=title]").Find("a").Text()
			d := s.Find("td[class=date]").Text()
			d1 := strings.Split(d, "(")[0]
			d2 := strings.Split(d, "(")[1]
			d2 = strings.Split(d2, ")")[1]
			n.Time, _ = time.ParseInLocation("2006/01/02 15:04", strings.TrimSpace(d1)+" "+strings.TrimSpace(d2), loc)
			nail = append(nail, n)
		}
	})
	return nail, nil
}

func showTitle() error {
	data, err := db.GetAllTitle()
	if err != nil {
		return err
	}
	sort.Sort(data)
	for _, d := range data {
		fmt.Println(d)
	}
	return nil
}
