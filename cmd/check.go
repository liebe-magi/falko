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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/reeve0930/foltia/db"
	"github.com/spf13/cobra"
)

type newAnimeInfo struct {
	TID     int
	Station string
	Title   string
	Time    time.Time
}

func (n newAnimeInfo) String() string {
	jweek := [7]string{"日", "月", "火", "水", "木", "金", "土"}
	return fmt.Sprintf("%s : %s(%d) %s",
		n.Time.Format("2006/01/02")+"("+jweek[n.Time.Weekday()]+") "+n.Time.Format("15:04"),
		n.Title,
		n.TID,
		n.Station,
	)
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the infomation of foltia",
	Run: func(cmd *cobra.Command, args []string) {
		newAnime, err := cmd.Flags().GetBool("new-anime")
		if err != nil {
			log.Fatalln(err)
		}
		tidFlag, err := cmd.Flags().GetBool("tid")
		if err != nil {
			log.Fatalln(err)
		}
		if checkFlag(newAnime, tidFlag) {
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
		} else {
			log.Println("Please check flag. You can use only one flag.")
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().BoolP("new-anime", "n", false, "Check the list of new anime")
	checkCmd.Flags().BoolP("tid", "t", false, "Check the list of anime title")
}

func checkFlag(n bool, t bool) bool {
	count := 0
	if n {
		count++
	}
	if t {
		count++
	}
	if count == 1 {
		return true
	}
	return false
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
