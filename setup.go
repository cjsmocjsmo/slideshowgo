package main

import (
	"image"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
)

type ImageData struct {
	Name string
	Path string
	Idx int
	Orientation string
}

func img_orient(imgPath string) (string, error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return "", err
	}

	if config.Width > config.Height {
		return "landscape", nil
	} else if config.Width < config.Height {
		return "portrait", nil
	} else {
		return "square", nil
	}
}

func create_img_db_table(dpath string) {
	dbPath := filepath.Join(dpath, "images.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("Failed to open database:", err)
		return
	}
	defer db.Close()

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS images (
		Name TEXT,
		Path TEXT,
		Idx INTEGER,
		Orientation TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		fmt.Println("Failed to create table:", err)
		return
	}
}

func Walk_Img_Dir(dbpath string, dir string) error {
	idx := 0

	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return err
	}
	defer db.Close()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if it's a regular file and has .jpg extension (case insensitive)
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".jpg") {
			idx += 1
			orientation, orientErr := img_orient(path)
			if orientErr != nil {
				return orientErr
			}
			imageData := ImageData{
				Name:        info.Name(),
				Path:        path,
				Idx:         idx,
				Orientation: orientation,
			}

			insertSQL := `INSERT INTO images (Name, Path, Idx, Orientation) VALUES (?, ?, ?, ?)`
			_, err = db.Exec(insertSQL, imageData.Name, imageData.Path, imageData.Idx, imageData.Orientation)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

	

