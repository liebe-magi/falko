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
	"path/filepath"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	homedir "github.com/mitchellh/go-homedir"
)

//NewAnime is a struct of video file
type NewAnime struct {
	gorm.Model
	TID     int
	Title   string
	Station string
	Time    time.Time
}

// InitNewAnimeDB : Initialize NewAnime DB
func InitNewAnimeDB() error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_newanime.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	db.AutoMigrate(&NewAnime{})
	return nil
}

// InsertNewAnime : Insert Data to NewAnime DB
func InsertNewAnime(tid int, title string, station string, time time.Time) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_newanime.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	db.Create(&NewAnime{TID: tid, Title: title, Station: station, Time: time})
	return nil
}

// UpdateNewAnime : Update Data of NewAnime DB
func UpdateNewAnime(id uint, tid int, title string, station string, time time.Time) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_newanime.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	var na NewAnime
	db.First(&na, id)
	na.TID = tid
	na.Title = title
	na.Station = station
	na.Time = time
	db.Save(&na)
	return nil
}

// DeleteNewAnime : Delete Data of NewAnime DB
func DeleteNewAnime(id uint) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_newanime.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	var vf NewAnime
	db.First(&vf, id)
	db.Delete(&vf)
	return nil
}

// GetAllNewAnime : Get All Data from NewAnime DB
func GetAllNewAnime() ([]NewAnime, error) {
	home, err := homedir.Dir()
	if err != nil {
		return []NewAnime{}, err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_newanime.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return []NewAnime{}, err
	}
	var nal []NewAnime
	db.Order("created_at desc").Find(&nal)
	return nal, nil
}

// GetOneNewAnime : Get Data from NewAnime DB
func GetOneNewAnime(id uint) (NewAnime, error) {
	home, err := homedir.Dir()
	if err != nil {
		return NewAnime{}, err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_newanime.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return NewAnime{}, err
	}
	var na NewAnime
	db.First(&na, id)
	return na, nil
}
