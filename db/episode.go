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

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	homedir "github.com/mitchellh/go-homedir"
)

//AnimeEpisode is a struct of anime episode
type AnimeEpisode struct {
	gorm.Model
	TID        int
	EpNum      int
	EpTitle   string
	CopyStatus bool
}

// InitEpisodeDB : Initialize Episode DB
func InitEpisodeDB() error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_episode.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	db.AutoMigrate(&AnimeEpisode{})
	return nil
}

// InsertEpisode : Insert data to Episode DB
func InsertEpisode(tid int, epnum int, eptitle string, copyStatus bool) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_episode.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	db.Create(&AnimeEpisode{TID: tid,EpNum: epnum, EpTitle: eptitle, CopyStatus: copyStatus})
	return nil
}

// UpdateEpisode : Update data of Episode DB
func UpdateEpisode(id uint, tid int, epnum int, eptitle string, copyStatus bool) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_episode.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	var ep AnimeEpisode
	db.First(&ep, id)
	ep.TID = tid
	ep.EpNum = epnum
	ep.EpTitle = eptitle
	ep.CopyStatus = copyStatus
	db.Save(&ep)
	return nil
}

// DeleteEpisode : Delete data of Episode DB
func DeleteEpisode(id uint) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_episode.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	var ep AnimeEpisode
	db.First(&ep, id)
	db.Delete(&ep)
	return nil
}

// GetAllEpisode : Get All Data from Episode DB
func GetAllEpisode() ([]AnimeEpisode, error) {
	home, err := homedir.Dir()
	if err != nil {
		return []AnimeEpisode{}, err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_episode.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return []AnimeEpisode{}, err
	}
	var epl []AnimeEpisode
	db.Order("created_at desc").Find(&epl)
	return epl, nil
}

// GetOneEpisode : Get Data from Episode DB
func GetOneEpisode(id uint) (AnimeEpisode, error) {
	home, err := homedir.Dir()
	if err != nil {
		return AnimeEpisode{}, err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_episode.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return AnimeEpisode{}, err
	}
	var ep AnimeEpisode
	db.First(&ep, id)
	return ep, nil
}
