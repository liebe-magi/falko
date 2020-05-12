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
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

type reservedProg struct {
	TID     int
	Title   string
	Station int
	Quality int
	MP2CM   int
	MP4CM   int
	Now     bool
}
type resevedProgList []reservedProg

type station struct {
	Name   string
	ID     int
	StType int
}
type stationList []station

func (r reservedProg) String() string {
	t := strconv.Itoa(r.TID)
	text := t + " : " + r.Title
	s, _ := getStation(r.Station)
	text += " " + s + " "
	text += fmt.Sprintf("[%d, %d, %d]", r.Quality, r.MP2CM, r.MP4CM)
	return text
}

// reserveCmd represents the reserve command
var reserveCmd = &cobra.Command{
	Use:   "reserve",
	Short: "TID指定による録画予約",
	Run: func(cmd *cobra.Command, args []string) {
		list, err := cmd.Flags().GetBool("list")
		if err != nil {
			log.Fatalln(err)
		}
		remove, err := cmd.Flags().GetBool("remove")
		if err != nil {
			log.Fatalln(err)
		}
		if list && remove {
			log.Fatalln(fmt.Errorf("2つのフラグを同時に指定することはできません"))
		} else if list && !remove {
			err = showReservedList()
			if err != nil {
				log.Println(err)
			}
		} else if !list && remove {
			err = dereserveProc(args)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			err = reserveProc(args)
			if err != nil {
				log.Fatalln(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(reserveCmd)

	reserveCmd.Flags().BoolP("list", "l", false, "録画予約している番組を表示")
	reserveCmd.Flags().BoolP("remove", "r", false, "予約の取消")
}

func reserveProc(args []string) error {
	if len(args) == 1 {
		tid, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		err = reserve(tid, 0, conf.encQuality, conf.mp2cut, conf.mp4cut)
	} else if len(args) == 2 {
		tid, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		sid, err := getStationID(args[1])
		if err != nil {
			return err
		}
		err = reserve(tid, sid, conf.encQuality, conf.mp2cut, conf.mp4cut)
	} else {
		return fmt.Errorf("引数の数が不正です")
	}
	return nil
}

func dereserveProc(args []string) error {
	if len(args) == 1 {
		tid, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		err = dereserve(tid, 0)
	} else if len(args) == 2 {
		tid, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		sid, err := getStationID(args[1])
		if err != nil {
			return err
		}
		err = dereserve(tid, sid)
	} else {
		return fmt.Errorf("引数の数が不正です")
	}
	return nil
}

func showReservedList() error {
	rpl, err := getReservedList()
	if err != nil {
		return err
	}
	fmt.Println("録画予約一覧")
	for _, r := range rpl {
		fmt.Println(r)
	}
	return nil
}

func getReservedList() (resevedProgList, error) {
	url := "http://" + conf.fHost + "/setup/listreserve.php"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return resevedProgList{}, err
	}
	var rpl resevedProgList
	doc.Find("#setUpTable > table > tbody > tr").Each(func(i int, s *goquery.Selection) {
		var rp reservedProg
		rp.TID, err = strconv.Atoi(s.Find("td:nth-child(2) > a").Text())
		if err != nil {
			log.Fatalln(err)
		}
		rp.Station, err = getStationID(strings.TrimSpace(s.Find("td:nth-child(3)").Text()))
		if err != nil {
			log.Fatalln(err)
		}
		rp.Title = strings.TrimSpace(s.Find("td:nth-child(4) > a").Text())
		rp.Quality, err = getQualityNum(strings.TrimSpace(s.Find("td:nth-child(5)").Text()))
		if err != nil {
			log.Fatalln(err)
		}
		rp.MP2CM, err = getCMEditNum(strings.TrimSpace(s.Find("td:nth-child(6)").Text()))
		if err != nil {
			log.Fatalln(err)
		}
		rp.MP4CM, err = getCMEditNum(strings.TrimSpace(s.Find("td:nth-child(7)").Text()))
		if err != nil {
			log.Fatalln(err)
		}
		pnum, err := strconv.Atoi(strings.TrimSpace(s.Find("td:nth-child(8)").Text()))
		if err != nil {
			log.Fatalln(err)
		}
		if pnum != 0 {
			rp.Now = true
		} else {
			rp.Now = false
		}
		rpl = append(rpl, rp)
	})
	return rpl, nil
}

func getStationList() stationList {
	var sl stationList
	sl = append(sl, station{Name: "NHK総合", ID: 1, StType: 0})
	sl = append(sl, station{Name: "NHK Eテレ", ID: 2, StType: 0})
	sl = append(sl, station{Name: "フジテレビ", ID: 3, StType: 0})
	sl = append(sl, station{Name: "日本テレビ", ID: 4, StType: 0})
	sl = append(sl, station{Name: "TBS", ID: 5, StType: 0})
	sl = append(sl, station{Name: "テレビ朝日", ID: 6, StType: 0})
	sl = append(sl, station{Name: "テレビ東京", ID: 7, StType: 0})
	sl = append(sl, station{Name: "tvk", ID: 8, StType: 0})
	sl = append(sl, station{Name: "NHK-BS1", ID: 9, StType: 0})
	sl = append(sl, station{Name: "NHK-BS2", ID: 10, StType: 0})
	sl = append(sl, station{Name: "チバテレビ", ID: 13, StType: 0})
	sl = append(sl, station{Name: "テレ玉", ID: 14, StType: 0})
	sl = append(sl, station{Name: "BSテレ東", ID: 15, StType: 1})
	sl = append(sl, station{Name: "BS-TBS", ID: 16, StType: 1})
	sl = append(sl, station{Name: "BSフジ", ID: 17, StType: 1})
	sl = append(sl, station{Name: "BS朝日", ID: 18, StType: 1})
	sl = append(sl, station{Name: "TOKYO MX", ID: 19, StType: 0})
	sl = append(sl, station{Name: "BS日テレ", ID: 71, StType: 1})
	sl = append(sl, station{Name: "BS11イレブン", ID: 128, StType: 1})
	sl = append(sl, station{Name: "BS12トゥエルビ", ID: 129, StType: 1})
	sl = append(sl, station{Name: "NHK BSプレミアム", ID: 179, StType: 1})
	sl = append(sl, station{Name: "TOKYO MX2", ID: 187, StType: 0})
	return sl
}

func getStation(id int) (string, error) {
	sl := getStationList()
	if id == 0 {
		return "[全局]", nil
	}
	for _, s := range sl {
		if id == s.ID {
			return s.Name, nil
		}
	}
	return "", fmt.Errorf("放送局IDが定義されていません : %d", id)
}

func getStationID(st string) (int, error) {
	sl := getStationList()
	if st == "[全局]" {
		return 0, nil
	}
	for _, s := range sl {
		if st == s.Name {
			return s.ID, nil
		}
	}
	return 0, fmt.Errorf("放送局名が定義されていません : %s", st)
}

func getQualityNum(s string) (int, error) {
	if s == "変換しない" {
		return 0, nil
	}
	if s == "SDのみ" {
		return 1, nil
	}
	if s == "HDのみ" {
		return 2, nil
	}
	if s == "SD+HD" {
		return 3, nil
	}
	return 0, fmt.Errorf("変換品質が定義されていません : %s", s)
}

func getQuality(v int) (string, error) {
	if v == 0 {
		return "変換しない", nil
	}
	if v == 1 {
		return "SDのみ", nil
	}
	if v == 2 {
		return "HDのみ", nil
	}
	if v == 3 {
		return "SD+HD", nil
	}
	return "", fmt.Errorf("変換品質の値が定義されていません : %d", v)
}

func getCMEditNum(s string) (int, error) {
	if s == "編集しない" {
		return 0, nil
	}
	if s == "本編のみ(CMカット)" {
		return 1, nil
	}
	if s == "CMのみ(本編カット)" {
		return 2, nil
	}
	if s == "本編+CM(同尺並び替え)" {
		return 3, nil
	}
	if s == "チャプタ追加" {
		return 4, nil
	}
	return 0, fmt.Errorf("CMカットルールが定義されていません : %s", s)
}

func getCMEdit(v int) (string, error) {
	if v == 0 {
		return "編集しない", nil
	}
	if v == 1 {
		return "本編のみ(CMカット)", nil
	}
	if v == 2 {
		return "CMのみ(本編カット)", nil
	}
	if v == 3 {
		return "本編+CM(同尺並び替え)", nil
	}
	if v == 4 {
		return "チャプタ追加", nil
	}

	return "", fmt.Errorf("CMカットルールの値が定義されていません : %d", v)
}

func reserve(tid int, station int, quality int, mp2cm int, mp4cm int) error {
	title, err := getTitle(tid)
	if err != nil {
		return err
	}
	s, err := getStation(station)
	if err != nil {
		return err
	}
	log.Printf("予約実行 : %s(%d) %s", title, tid, s)
	url := "http://" + conf.fHost + "/reservation/reservecomp.php"
	url += fmt.Sprintf("?station=%d", station)
	url += fmt.Sprintf("&transcodequality=%d", quality)
	url += fmt.Sprintf("&cmeditrulempeg2=%d", mp2cm)
	url += fmt.Sprintf("&cmeditrulemp4=%d", mp4cm)
	url += fmt.Sprintf("&usedigital=1&tid=%d", tid)
	_, err = goquery.NewDocument(url)
	return err
}

func dereserve(tid int, station int) error {
	title, err := getTitle(tid)
	if err != nil {
		return err
	}
	s, err := getStation(station)
	if err != nil {
		return err
	}
	log.Printf("予約取消 : %s(%d) %s", title, tid, s)
	url := "http://" + conf.fHost + "/reservation/delreserve.php"
	url += fmt.Sprintf("?sid=%d", station)
	url += fmt.Sprintf("&delflag=1&tid=%d", tid)
	_, err = goquery.NewDocument(url)
	return err
}
