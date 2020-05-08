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

package db

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	homedir "github.com/mitchellh/go-homedir"
)

//VideoFile is a struct of video file
type VideoFile struct {
	gorm.Model
	TID       int
	EpNum     int
	PID       int
	FileTS    string
	FileMP4HD string
	FileMP4SD string
	Station   string
	Time      time.Time
	Drop      int
	Scramble  int
}

func (v VideoFile) String() string {
	return fmt.Sprintf("PID : %d, Time : %s, Staiton : %s, FileTS : %s, FileMP4HD : %s, FileMP4SD : %s, Drop : %d, Scramble : %d\n", v.PID, v.Time, v.Station, v.FileTS, v.FileMP4HD, v.FileMP4SD, v.Drop, v.Scramble)
}

// InitVideoFileDB : Initialize VideoFile DB
func InitVideoFileDB() error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_videofile.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	db.AutoMigrate(&VideoFile{})
	return nil
}

// InsertVideoFile : Insert Data to VideoFile DB
func InsertVideoFile(tid int, epnum int, pid int, filets string, filemp4hd string, filemp4sd string, station string, time time.Time, drop int, scramble int) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_videofile.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	db.Create(&VideoFile{TID: tid, EpNum: epnum, PID: pid, FileTS: filets, FileMP4HD: filemp4hd, FileMP4SD: filemp4sd, Station: station, Time: time, Drop: drop, Scramble: scramble})
	return nil
}

// UpdateVideoFile : Update Data of VideoFile DB
func UpdateVideoFile(id uint, tid int, epnum int, pid int, filets string, filemp4hd string, filemp4sd string, station string, time time.Time, drop int, scramble int) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_videofile.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	var vf VideoFile
	db.First(&vf, id)
	vf.TID = tid
	vf.EpNum = epnum
	vf.PID = pid
	vf.FileTS = filets
	vf.FileMP4HD = filemp4hd
	vf.FileMP4SD = filemp4sd
	vf.Station = station
	vf.Time = time
	vf.Drop = drop
	vf.Scramble = scramble
	db.Save(&vf)
	return nil
}

// DeleteVideoFile : Delete Data of VideoFile DB
func DeleteVideoFile(id uint) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_videofile.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	var vf VideoFile
	db.First(&vf, id)
	db.Delete(&vf)
	return nil
}

// GetAllVideoFile : Get All Data from VideoFile DB
func GetAllVideoFile() ([]VideoFile, error) {
	home, err := homedir.Dir()
	if err != nil {
		return []VideoFile{}, err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_videofile.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return []VideoFile{}, err
	}
	var vfl []VideoFile
	db.Order("created_at desc").Find(&vfl)
	return vfl, nil
}

// GetOneVideoFile : Get Data from VideoFile DB
func GetOneVideoFile(id uint) (VideoFile, error) {
	home, err := homedir.Dir()
	if err != nil {
		return VideoFile{}, err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_videofile.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return VideoFile{}, err
	}
	var vf VideoFile
	db.First(&vf, id)
	return vf, nil
}
