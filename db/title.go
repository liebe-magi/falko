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

//AnimeTitle is a struct of anime title
type AnimeTitle struct {
	gorm.Model
	TID   int
	Title string
}

// InitTitleDB : Initialize Title DB
func InitTitleDB() error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_title.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	db.AutoMigrate(&AnimeTitle{})
	return nil
}

// InsertTitle : Insert Data to Title DB
func InsertTitle(tid int, title string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_title.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	db.Create(&AnimeTitle{TID: tid, Title: title})
	return nil
}

// UpdateTitle : Update Data of Title DB
func UpdateTitle(id uint, tid int, title string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_title.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	var t AnimeTitle
	db.First(&t, id)
	t.TID = tid
	t.Title = title
	db.Save(&t)
	return nil
}

// DeleteTitle : Delete Data of Title DB
func DeleteTitle(id uint) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_title.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return err
	}
	var t AnimeTitle
	db.First(&t, id)
	db.Delete(&t)
	return nil
}

// GetAllTitle : Get All Data from Title DB
func GetAllTitle() ([]AnimeTitle, error) {
	home, err := homedir.Dir()
	if err != nil {
		return []AnimeTitle{}, err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_title.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return []AnimeTitle{}, err
	}
	var tl []AnimeTitle
	db.Order("created_at desc").Find(&tl)
	return tl, nil
}

// GetOneTitle : Get Data from Title DB
func GetOneTitle(id uint) (AnimeTitle, error) {
	home, err := homedir.Dir()
	if err != nil {
		return AnimeTitle{}, err
	}
	dbPath := filepath.Join(home, ".config", "foltia", "foltia_title.sqlite3")
	db, err := gorm.Open("sqlite3", dbPath)
	defer db.Close()
	if err != nil {
		return AnimeTitle{}, err
	}
	var t AnimeTitle
	db.First(&t, id)
	return t, nil
}
