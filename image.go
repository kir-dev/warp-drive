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
	Title            string
	OriginalFilename string
	Height           int
	Width            int
	Hash             string
	Created          time.Time

	content []byte
}

func (img *imageRecord) path() string {
	return fmt.Sprintf("%s/%s/%s%s",
		img.Hash[0:2],
		img.Hash[0:8],
		img.Hash,
		filepath.Ext(img.OriginalFilename))
}

func (img *imageRecord) save() error {
	err1 := img.saveToDb()
	err2 := img.saveToDisk()

	if err1 != nil || err2 != nil {
		return fmt.Errorf("Could not save to disk (%v) or to the database(%v),", err2, err1)
	}
	log.Printf("Saved image. Filename: %s, hash: %s", img.OriginalFilename, img.Hash)
	return nil
}

func (img *imageRecord) saveToDb() error {
	_, err := imageInsertStmt.Exec(img.Title, img.OriginalFilename, img.path(), img.Height, img.Width, img.Hash, img.Created)
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

		if out.Chmod(0644) != nil {
			return "", err
		}

		err = jpeg.Encode(out, newImg, nil)
		if err != nil {
			return "", err
		}
		log.Printf("Resized image: %s", filename)
	}

	return newpath, nil
}
