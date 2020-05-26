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
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cheggaaa/pb/v3"
	"github.com/reeve0930/falko/db"
	"github.com/spf13/cobra"
)

type animeTitleInfo struct {
	TID   int
	Title string
	Yomi  string
	Year  int
}

type animeFileInfo struct {
	TID       int
	Title     string
	EpNum     int
	PID       int
	EpTitle   string
	Time      time.Time
	Station   string
	FileTS    string
	FileMP4HD string
	FileMP4SD string
}

var barTemp = `{{counters .}} {{bar . "|" "=" ">" "_" "|"}} {{ speed .}} {{percent .}} {{rtime . "ETA %s"}}`

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update local DB",
	Run: func(cmd *cobra.Command, args []string) {
		updateDB()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func updateDB() {
	log.Println("ローカルDBの更新を開始")

	db.InitTitleDB()
	db.InitEpisodeDB()
	db.InitVideoFileDB()

	log.Println("アニメタイトルDBを更新")
	atil, err := getAnimeTitleInfo()
	if err != nil {
		log.Fatalln(err)
	}
	err = insertNewTitle(atil)
	if err != nil {
		log.Fatalln(err)
	}
	err = activateAnimeTitle()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("動画ファイルの情報取得を開始")
	afil, err := getVideoFile()
	if err != nil {
		log.Fatalln(err)
	}
	err = removeVideoFile(afil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("エピソードDBの更新を開始")
	err = insertNewEpisode(afil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("動画ファイルDBの更新を開始")
	err = insertNewVideoFile(afil)
	if err != nil {
		log.Fatalln(err)
	}
	data, err := getCopyList()
	log.Printf("%d個の動画ファイルを検出", len(data))
	log.Println("ローカルDBの更新を完了")
}

func getAnimeTitleInfo() ([]animeTitleInfo, error) {
	url := "http://cal.syoboi.jp/db.php?Command=TitleLookup&TID=*"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return []animeTitleInfo{}, err
	}
	var atil []animeTitleInfo
	if err != nil {
		return []animeTitleInfo{}, err
	}
	doc.Find("TitleLookupResponse > TitleItems > TitleItem").Each(func(i int, s *goquery.Selection) {
		var a animeTitleInfo
		a.TID, _ = strconv.Atoi(s.Find("TID").Text())
		a.Title = strings.TrimSpace(s.Find("Title").Text())
		a.Yomi = strings.TrimSpace(s.Find("TitleYomi").Text())
		a.Year, _ = strconv.Atoi(s.Find("FirstYear").Text())
		atil = append(atil, a)
	})
	return atil, nil
}

func activateAnimeTitle() error {
	url := "http://" + conf.fHost + "/recorded/recfiles_tid.php?mode=detail"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}
	data, err := db.GetAllTitle()
	if err != nil {
		return err
	}
	num := doc.Find("#Librarytable > table > tbody").Find("tr").Length()
	bar := pb.ProgressBarTemplate(barTemp).Start(len(data))
	for _, d := range data {
		exists := false
		for i := 0; i < num; i++ {
			t := doc.Find(fmt.Sprintf("#Librarytable > table > tbody > tr:nth-child(%d) > td:nth-child(1) > a", i+1)).Text()
			tid, err := strconv.Atoi(t)
			if err != nil {
				log.Fatalln(err)
			}
			if tid == d.TID {
				if !d.Active {
					db.UpdateTitle(d.ID, d.TID, d.Title, d.TitleYomi, d.Year, true)
				}
				exists = true
				break
			}
		}
		if !exists && d.Active {
			db.UpdateTitle(d.ID, d.TID, d.Title, d.TitleYomi, d.Year, false)
		}
		bar.Increment()
	}
	bar.Finish()
	return nil
}

func insertNewTitle(atil []animeTitleInfo) error {
	data, err := db.GetAllTitle()
	if err != nil {
		return err
	}
	for _, a := range atil {
		exist := false
		for _, d := range data {
			if a.TID == d.TID {
				exist = true
				break
			}
		}
		if !exist {
			db.InsertTitle(a.TID, a.Title, a.Yomi, a.Year, false)
		}
	}
	return nil
}

func removeVideoFile(afil []animeFileInfo) error {
	data, err := db.GetAllVideoFile()
	if err != nil {
		return err
	}
	for _, d := range data {
		exists := false
		for _, a := range afil {
			if d.PID == a.PID {
				exists = true
				if d.FileTS != a.FileTS || d.FileMP4HD != a.FileMP4HD || d.FileMP4SD != a.FileMP4SD {
					log.Printf("動画ファイルの情報を更新 : %s (%d:%s)", a.Title, a.EpNum, a.EpTitle)
					err = db.UpdateVideoFile(d.ID, d.TID, d.EpNum, d.PID, a.FileTS, a.FileMP4HD, a.FileMP4SD, d.Station, d.Time, d.Drop, d.Scramble)
					if err != nil {
						return err
					}
				}
				break
			}
		}
		if !exists {
			title, err := getTitle(d.TID)
			if err != nil {
				return err
			}
			log.Printf("動画ファイルの情報を削除 : %s (%d): %d", title, d.EpNum, d.PID)
			err = db.DeleteVideoFile(d.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getTitle(tid int) (string, error) {
	data, err := db.GetAllTitle()
	if err != nil {
		return "", err
	}
	for _, d := range data {
		if d.TID == tid {
			return d.Title, nil
		}
	}
	return "", fmt.Errorf("TIDが未定義")
}

func getVideoFile() ([]animeFileInfo, error) {
	var afil []animeFileInfo
	data, err := db.GetAllTitle()
	if err != nil {
		return []animeFileInfo{}, err
	}
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return []animeFileInfo{}, err
	}
	bar := pb.ProgressBarTemplate(barTemp).Start(len(data))
	for _, d := range data {
		if d.Active {
			url := "http://" + conf.fHost + "/recorded/recfiles_tid.php?mode=detail&tid=" + fmt.Sprintf("%d", d.TID)
			doc, err := goquery.NewDocument(url)
			if err != nil {
				continue
			}
			num := doc.Find("#libraryDetail > li").Length()
			for i := 0; i < num; i++ {
				var afi animeFileInfo
				afi.TID = d.TID
				afi.Title = d.Title
				e := strings.TrimPrefix(doc.Find(fmt.Sprintf("#libraryDetail > li:nth-child(%d) > div.programInfo > ul > li:nth-child(2)", i+1)).Text(), "話数：")
				if e == "[話数]" {
					e = "-1"
				}
				afi.EpNum, err = strconv.Atoi(e)
				if err != nil {
					return []animeFileInfo{}, err
				}
				afi.EpTitle = strings.TrimSpace(strings.TrimPrefix(doc.Find(fmt.Sprintf("#libraryDetail > li:nth-child(%d) > div.programInfo > ul > li:nth-child(3)", i+1)).Text(), "サブタイトル："))
				t := strings.TrimSpace(strings.TrimPrefix(doc.Find(fmt.Sprintf("#libraryDetail > li:nth-child(%d) > div.programInfo > ul > li:nth-child(4)", i+1)).Text(), "録画日時："))
				afi.Time, err = time.ParseInLocation("2006/01/02 15:04", strings.Split(t, "(")[0]+strings.Split(t, ")")[1], loc)
				if err != nil {
					return []animeFileInfo{}, err
				}
				afi.Station = strings.TrimSpace(strings.TrimPrefix(doc.Find(fmt.Sprintf("#libraryDetail > li:nth-child(%d) > div.programInfo > ul > li:nth-child(5)", i+1)).Text(), "放送局："))
				status := strings.TrimSpace(strings.TrimPrefix(doc.Find(fmt.Sprintf("#libraryDetail > li:nth-child(%d) > div.programInfo > ul > li:nth-child(6)", i+1)).Text(), "ステータス："))
				if status != "完了" {
					continue
				}
				doc.Find(fmt.Sprintf("#libraryDetail > li:nth-child(%d) > div.programInfo > div > ul.fileType > li", i+1)).Each(func(j int, s *goquery.Selection) {
					t, _ := s.Attr("class")
					if t == "mpeg2" {
						afi.FileTS = strings.TrimSpace(s.Text())
					} else if t == "mp4HD" {
						afi.FileMP4HD = strings.TrimSpace(s.Text())
					} else if t == "mp4SD" {
						afi.FileMP4SD = strings.TrimSpace(s.Text())
					}
				})
				p, exists := doc.Find(fmt.Sprintf("#libraryDetail > li:nth-child(%d) > div.programInfo > ul > div > a ", i+1)).Attr("href")
				if !exists {
					return []animeFileInfo{}, fmt.Errorf("PID not found")
				}
				afi.PID, err = strconv.Atoi(strings.TrimPrefix(p, "./selectcaptureimage.php?pid="))
				if err != nil {
					return []animeFileInfo{}, err
				}
				afil = append(afil, afi)
			}
		}
		bar.Increment()
	}
	bar.Finish()
	return afil, nil
}

func checkEpisode(tid int, epnum int) (bool, error) {
	data, err := db.GetAllEpisode()
	if err != nil {
		return false, err
	}
	for _, d := range data {
		if d.TID == tid && d.EpNum == epnum {
			return true, nil
		}
	}
	return false, nil
}

func checkVideoFile(pid int) (bool, error) {
	data, err := db.GetAllVideoFile()
	if err != nil {
		return false, err
	}
	for _, d := range data {
		if d.PID == pid {
			return true, nil
		}
	}
	return false, nil
}

func insertNewEpisode(afil []animeFileInfo) error {
	bar := pb.ProgressBarTemplate(barTemp).Start(len(afil))
	for _, a := range afil {
		exists, err := checkEpisode(a.TID, a.EpNum)
		if err != nil {
			return err
		}
		if !exists {
			err = db.InsertEpisode(a.TID, a.EpNum, a.EpTitle, false)
			if err != nil {
				return err
			}
		}
		bar.Increment()
	}
	bar.Finish()
	return nil
}

func insertNewVideoFile(afil []animeFileInfo) error {
	bar := pb.ProgressBarTemplate(barTemp).Start(len(afil))
	for _, a := range afil {
		exists, err := checkVideoFile(a.PID)
		if err != nil {
			return err
		}
		if !exists {
			dr, sc, err := getTSInfo(a.PID)
			if err != nil {
				return err
			}
			err = db.InsertVideoFile(a.TID, a.EpNum, a.PID, a.FileTS, a.FileMP4HD, a.FileMP4SD, a.Station, a.Time, dr, sc)
			if err != nil {
				return err
			}
		}
		bar.Increment()
	}
	bar.Finish()
	return nil
}

func getTSInfo(pid int) (int, int, error) {
	url := "http://" + conf.fHost + "/recorded/showcminfo.php?pid=" + fmt.Sprintf("%d", pid)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return 0, 0, err
	}
	n := doc.Find("#programInfo").Children().Length()
	text := strings.Split(doc.Find(fmt.Sprintf("#programInfo > p:nth-child(%d)", n)).Text(), "\n")
	dr := 0
	sc := 0
	for i, t := range text {
		if i > 0 {
			e := strings.Split(t, ", ")
			if len(e) == 4 {
				d, err := strconv.Atoi(strings.TrimSpace(strings.Split(e[2], "=")[1]))
				if err != nil {
					return 0, 0, err
				}
				s, err := strconv.Atoi(strings.TrimSpace(strings.Split(e[3], "=")[1]))
				if err != nil {
					return 0, 0, err
				}
				dr += d
				sc += s
			}
		}
	}
	return dr, sc, nil
}
