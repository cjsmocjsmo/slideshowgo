package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"
	_ "github.com/mattn/go-sqlite3"
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
		fmt.Println("Landscape")
		return "landscape", nil
	} else if config.Width < config.Height {
		fmt.Println("Portrait")
		return "portrait", nil
	} else {
		fmt.Println("Square")
		return "square", nil
	}
}

func calc_name(imgpath string) string {
	file, err := os.Open(imgpath)
	if err != nil {
		return ""
	}
	defer file.Close()

	hasher := md5.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", hasher.Sum(nil))
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
				Name:        calc_name(path),
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

	

