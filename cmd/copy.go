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
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/reeve0930/falko/db"
	"github.com/spf13/cobra"
)

type fileCopyInfo struct {
	tid      int
	title    string
	epNum    int
	epTitle  string
	pid      int
	srcname  string
	dstname  string
	station  string
	scramble bool
}

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy (TID) (Episode No)",
	Short: "動画ファイルのコピー",
	Run: func(cmd *cobra.Command, args []string) {
		list, err := cmd.Flags().GetBool("list")
		if err != nil {
			log.Fatalln(err)
		}
		reset, err := cmd.Flags().GetBool("reset")
		if err != nil {
			log.Fatalln(err)
		}
		ignore, err := cmd.Flags().GetBool("ignoreDrop")
		if err != nil {
			log.Fatalln(err)
		}
		tid := -1
		epNum := -1
		if len(args) == 1 {
			tid, err = strconv.Atoi(args[0])
			if err != nil {
				log.Fatalln(err)
			}
		} else if len(args) == 2 {
			tid, err = strconv.Atoi(args[0])
			if err != nil {
				log.Fatalln(err)
			}
			epNum, err = strconv.Atoi(args[1])
			if err != nil {
				log.Fatalln(err)
			}
		}
		if list && !reset {
			err = showCopyList(tid, epNum, ignore)
			if err != nil {
				log.Fatalln(err)
			}
		} else if !list && reset && !ignore {
			if len(args) != 2 {
				log.Fatalln("TIDとエピソード番号を指定して下さい")
			} else {
				err = resetCopyStatus(tid, epNum)
				if err != nil {
					log.Fatalln(err)
				}
			}
		} else if !list && !reset {
			err = copyFiles(tid, epNum, ignore)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Fatalln("フラグの指定を確認して下さい")
		}
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	copyCmd.Flags().BoolP("list", "l", false, "コピー予定のファイル一覧を表示")
	copyCmd.Flags().BoolP("reset", "r", false, "動画ファイルのコピー済みフラグを削除")
	copyCmd.Flags().BoolP("ignoreDrop", "i", false, "TSドロップを無視してコピー")
}

func showCopyList(tid int, epNum int, ignore bool) error {
	fcil, err := getCopyList(ignore)
	if err != nil {
		return err
	}
	if tid != -1 {
		fcil, err = filter(fcil, tid, epNum)
		if err != nil {
			return err
		}
	}
	showList(fcil)
	log.Printf("%d個の動画ファイルを検出", len(fcil))
	return nil
}

func filter(fcil []fileCopyInfo, tid int, epNum int) ([]fileCopyInfo, error) {
	var fcilNew []fileCopyInfo
	for _, f := range fcil {
		if epNum == -1 {
			if f.tid == tid {
				fcilNew = append(fcilNew, f)
			}
		} else {
			if f.tid == tid && f.epNum == epNum {
				fcilNew = append(fcilNew, f)
			}
		}
	}
	return fcilNew, nil
}

func copyFiles(tid int, epNum int, ignore bool) error {
	log.Println("コピー開始")
	fcil, err := getCopyList(ignore)
	if err != nil {
		return err
	}
	if tid != -1 {
		fcil, err = filter(fcil, tid, epNum)
		if err != nil {
			return err
		}
	}
	ep, err := db.GetAllEpisode()
	if err != nil {
		return err
	}
	key, err := db.GetAllKeywordRecFile()
	if err != nil {
		return err
	}
	for i, f := range fcil {
		log.Printf("[%d/%d] %s (%d:%s)", i+1, len(fcil), f.title, f.epNum, f.epTitle)
		if f.scramble {
			log.Println("スクランブルが未解除")
			f.dstname = "[S]" + f.dstname
		}
		src := filepath.Join(conf.fPath, f.srcname)
		f.dstname = fixFileNameLength(f.dstname)
		dst := filepath.Join(conf.cDest, f.dstname)
		err = copyVideoFile(src, dst)
		if err != nil {
			return err
		}
		if f.tid != -1 {
			for _, e := range ep {
				if !ignore && (e.TID == f.tid && e.EpNum == f.epNum) {
					db.UpdateEpisode(e.ID, e.TID, e.EpNum, e.EpTitle, true)
					break
				}
			}
		} else {
			for _, k := range key {
				if f.pid == k.PID {
					db.UpdateKeywordRecFile(k.ID, k.Keyword, k.Title, k.PID, k.FileTS, k.FileMP4HD, k.FileMP4SD, k.Station, k.Time, k.Drop, k.Scramble, true)
					break
				}
			}
		}
	}
	log.Println("コピー完了")
	return nil
}

func fixFileNameLength(name string) string {
	if len(name) <= 255 {
		return name
	}
	n := strings.Split(name, ".")
	s := []rune(n[0])
	nn := fmt.Sprintf(string(s[:(len(s) - 1)]))
	return fixFileNameLength(nn + "." + n[1])
}

func copyVideoFile(src string, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	sourceStat, err := s.Stat()
	if err != nil {
		return err
	}
	srcSize := sourceStat.Size()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	bar := pb.New64(srcSize).SetTemplateString(barTemp).Start()
	reader := bar.NewProxyReader(s)

	for {
		_, err = io.Copy(d, reader)
		if err != nil {
			log.Println(err)
			log.Println("コピー処理が失敗しました。リトライします。")
		} else {
			break
		}
	}
	bar.Finish()

	return nil
}

func getCopyList(ignore bool) ([]fileCopyInfo, error) {
	title, err := db.GetAllTitle()
	if err != nil {
		return []fileCopyInfo{}, err
	}
	episode, err := db.GetAllEpisode()
	if err != nil {
		return []fileCopyInfo{}, err
	}
	videofile, err := db.GetAllVideoFile()
	if err != nil {
		return []fileCopyInfo{}, err
	}
	var fcil []fileCopyInfo
	for _, t := range title {
		for _, e := range episode {
			if e.CopyStatus {
				continue
			}
			if t.TID == e.TID {
				var f fileCopyInfo
				f.tid = t.TID
				f.title = t.Title
				f.epNum = e.EpNum
				f.epTitle = e.EpTitle
				f.dstname = conf.cFilename
				f.dstname = strings.Replace(f.dstname, "%title%", t.Title, -1)
				f.dstname = strings.Replace(f.dstname, "%epnum%", fmt.Sprintf("%02d", e.EpNum), -1)
				f.dstname = strings.Replace(f.dstname, "%eptitle%", e.EpTitle, -1)
				if conf.cFiletype == "TS" {
					f.dstname = f.dstname + ".ts"
				} else if conf.cFiletype == "MP4" {
					f.dstname = f.dstname + ".mp4"
				} else {
					return []fileCopyInfo{}, fmt.Errorf("設定が異常値 : copy_filetype")
				}

				nonDropExists := false
				fileExists := false
				for _, v := range videofile {
					if e.TID == v.TID && e.EpNum == v.EpNum {
						fileExists = true
						if ignore || (v.Drop < conf.cDropThresh) {
							nonDropExists = true
							if f.srcname == "" {
								name, check := getSrcname(v)
								if check {
									f.srcname = name
									f.pid = v.PID
									f.station = v.Station
									if v.Scramble != 0 {
										f.scramble = true
									} else {
										f.scramble = false
									}
								} else {
									continue
								}
							} else {
								p1, err := getStationPriority(f.station)
								if err != nil {
									return []fileCopyInfo{}, err
								}
								p2, err := getStationPriority(v.Station)
								if err != nil {
									return []fileCopyInfo{}, err
								}
								if p2 > p1 {
									name, check := getSrcname(v)
									if check {
										f.srcname = name
										f.pid = v.PID
										f.station = v.Station
										if v.Scramble != 0 {
											f.scramble = true
										} else {
											f.scramble = false
										}
									}
								}
							}
						}
					}
				}
				if f.srcname != "" {
					fcil = append(fcil, f)
				} else {
					if !nonDropExists && fileExists {
						log.Printf("設定値を超えたTSドロップが発生 : %s (%d:%s)", t.Title, e.EpNum, e.EpTitle)
					}
				}
			}
		}
	}
	//キーワード録画を追加
	key, err := db.GetAllKeywordRecFile()
	if err != nil {
		return []fileCopyInfo{}, err
	}
	for _, k := range key {
		if !k.Copy {
			var fci fileCopyInfo
			fci.title = k.Title
			fci.tid = -1
			fci.epNum = -1
			fci.epTitle = ""
			fci.station = k.Station
			fci.pid = k.PID
			if k.Scramble != 0 {
				fci.scramble = true
			} else {
				fci.scramble = false
			}
			name, check := getSrcnameKey(k)
			if check {
				fci.srcname = name
				d := k.Time.Format("20060102150405")
				t := strings.Replace(k.Title, "?", "？", -1)
				t = strings.Replace(k.Title, "\"", "”", -1)
				fci.dstname = fmt.Sprintf("[D%d]%s(%s)_%s_%s", k.Drop, k.Keyword, k.Station, d, t)
				if conf.cFiletype == "TS" {
					fci.dstname = fci.dstname + ".ts"
				} else if conf.cFiletype == "MP4" {
					fci.dstname = fci.dstname + ".mp4"
				} else {
					return []fileCopyInfo{}, fmt.Errorf("設定が異常値 : copy_filetype")
				}
				fcil = append(fcil, fci)
			}
		}
	}
	return fcil, nil
}

func getStationPriority(st string) (int, error) {
	station := getStationList()
	for _, s := range station {
		if s.Name == st {
			return s.StType, nil
		}
	}
	return -1, fmt.Errorf("放送局名が未定義 : %s", st)
}

func getSrcname(v db.VideoFile) (string, bool) {
	if conf.cFiletype == "TS" {
		if v.FileTS != "" {
			return v.FileTS, true
		}
		return "", false
	}
	if v.FileMP4HD != "" {
		return v.FileMP4HD, true
	}
	if v.FileMP4SD != "" {
		return v.FileMP4SD, true
	}
	return "", false
}

func getSrcnameKey(k db.KeywordRecFile) (string, bool) {
	if conf.cFiletype == "TS" {
		if k.FileTS != "" {
			return k.FileTS, true
		}
		return "", false
	}
	if k.FileMP4HD != "" {
		return k.FileMP4HD, true
	}
	if k.FileMP4SD != "" {
		return k.FileMP4SD, true
	}
	return "", false
}

func showList(fcil []fileCopyInfo) {
	for _, f := range fcil {
		fmt.Printf("%d : %s\n", f.pid, f.dstname)
	}
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
			log.Printf("コピー済みフラグをリセット : (%d)%s (%d:%s)", d.TID, title, d.EpNum, d.EpTitle)
			err = db.UpdateEpisode(d.ID, d.TID, d.EpNum, d.EpTitle, false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
