package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	ImageInsertSql  = "INSERT INTO images(title, orig_filename, filepath, height, width, hash, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	ImageGetPathSql = "SELECT filepath, width FROM images WHERE hash = $1"
)

type imageRecord struct {
	title            string
	originalFilename string
	height           int
	width            int
	hash             string
	created          time.Time

	content []byte
}

func (img *imageRecord) path() string {
	return fmt.Sprintf("%s/%s/%s%s",
		img.hash[0:2],
		img.hash[0:8],
		img.hash,
		filepath.Ext(img.originalFilename))
}

func (img *imageRecord) save() error {
	err1 := img.saveToDb()
	err2 := img.saveToDisk()

	if err1 != nil || err2 != nil {
		return fmt.Errorf("Could not save to disk (%v) or to the database(%v),", err2, err1)
	}
	log.Printf("Saved image. Filename: %s, hash: %s", img.originalFilename, img.hash)
	return nil
}

func (img *imageRecord) saveToDb() error {
	_, err := imageInsertStmt.Exec(img.title, img.originalFilename, img.path(), img.height, img.width, img.hash, img.created)
	return err
}

func (img *imageRecord) saveToDisk() error {
	fullPath := config.absolutePath(img.path())
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(fullPath, img.content, 0644)
}

func getImagePath(hash string, width int) (string, error) {
	var path string
	var actualWidth int
	if err := getImagePathStmt.QueryRow(hash).Scan(&path, &actualWidth); err != nil {
		return "", err
	}

	if filepath.Ext(path) == ".gif" {
		return "", fmt.Errorf("resizing gifs are not supported")
	}

	if width <= 0 || width >= actualWidth {
		return config.absolutePath(path), nil
	}

	return resizeImage(config.absolutePath(path), width)
}

// Resize the image if neccessary and save it. Returns the new path of the
// resized image.
func resizeImage(path string, width int) (string, error) {
	dir := filepath.Dir(path)
	filename := fmt.Sprintf("%d_%s", width, filepath.Base(path))
	newpath := filepath.Join(dir, filename)

	if _, err := os.Stat(newpath); os.IsNotExist(err) {
		file, err := os.Open(path)
		if err != nil {
			return "", err
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return "", err
		}

		newImg := resize.Resize(uint(width), 0, img, resize.Lanczos3)

		out, err := os.Create(newpath)
		if err != nil {
			return "", err
		}
		defer out.Close()

		err = jpeg.Encode(out, newImg, nil)
		if err != nil {
			return "", err
		}
		log.Printf("Resized image: %s", filename)
	}

	return newpath, nil
}
