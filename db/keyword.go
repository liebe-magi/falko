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

package db

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	homedir "github.com/mitchellh/go-homedir"
)

// KeywordRecFile is a struct of video file
type KeywordRecFile struct {
	gorm.Model
	Keyword   string
	Title     string
	PID       int
	FileTS    string
	FileMP4HD string
	FileMP4SD string
	Station   string
	Time      time.Time
	Drop      int
	Scramble  int
	Copy      bool
}

func (v KeywordRecFile) String() string {
	return fmt.Sprintf("Keyword: %s, Title: %s, PID : %d, Time : %s, Staiton : %s, FileTS : %s, FileMP4HD : %s, FileMP4SD : %s, Drop : %d, Scramble : %d, Copy : %t\n", v.Keyword, v.Title, v.PID, v.Time, v.Station, v.FileTS, v.FileMP4HD, v.FileMP4SD, v.Drop, v.Scramble, v.Copy)
}

// InitKeywordRecFileDB : Initialize KeywordRecFile DB
func InitKeywordRecFileDB() error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "falko", "foltia_keyword.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	db.AutoMigrate(&KeywordRecFile{})
	return nil
}

// InsertKeywordRecFile : Insert Data to KeywordRecFile DB
func InsertKeywordRecFile(keyword string, title string, pid int, filets string, filemp4hd string, filemp4sd string, station string, time time.Time, drop int, scramble int, cp bool) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "falko", "foltia_keyword.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	db.Create(&KeywordRecFile{Keyword: keyword, Title: title, PID: pid, FileTS: filets, FileMP4HD: filemp4hd, FileMP4SD: filemp4sd, Station: station, Time: time, Drop: drop, Scramble: scramble, Copy: cp})
	return nil
}

// UpdateKeywordRecFile : Update Data of KeywordRecFile DB
func UpdateKeywordRecFile(id uint, keyword string, title string, pid int, filets string, filemp4hd string, filemp4sd string, station string, time time.Time, drop int, scramble int, cp bool) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "falko", "foltia_keyword.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	var krf KeywordRecFile
	db.First(&krf, id)
	krf.Keyword = keyword
	krf.Title = title
	krf.PID = pid
	krf.FileTS = filets
	krf.FileMP4HD = filemp4hd
	krf.FileMP4SD = filemp4sd
	krf.Station = station
	krf.Time = time
	krf.Drop = drop
	krf.Scramble = scramble
	krf.Copy = cp
	db.Save(&krf)
	return nil
}

// DeleteKeywordRecFile : Delete Data of KeywordRecFile DB
func DeleteKeywordRecFile(id uint) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "falko", "foltia_keyword.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	var krf KeywordRecFile
	db.First(&krf, id)
	db.Delete(&krf)
	return nil
}

// GetAllKeywordRecFile : Get All Data from KeywordRecFile DB
func GetAllKeywordRecFile() ([]KeywordRecFile, error) {
	home, err := homedir.Dir()
	if err != nil {
		return []KeywordRecFile{}, err
	}
	dbPath := filepath.Join(home, ".config", "falko", "foltia_keyword.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return []KeywordRecFile{}, err
	}
	defer db.Close()
	var krfl []KeywordRecFile
	db.Order("created_at desc").Find(&krfl)
	return krfl, nil
}

// GetOneKeywordRecFile : Get Data from KeywordRecFile DB
func GetOneKeywordRecFile(id uint) (KeywordRecFile, error) {
	home, err := homedir.Dir()
	if err != nil {
		return KeywordRecFile{}, err
	}
	dbPath := filepath.Join(home, ".config", "falko", "foltia_keyword.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return KeywordRecFile{}, err
	}
	defer db.Close()
	var krf KeywordRecFile
	db.First(&krf, id)
	return krf, nil
}
