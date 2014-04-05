package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	ImageInsertSql = "INSERT INTO images(title, orig_filename, filepath, height, width, hash, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
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
	return fmt.Sprintf("%s/%s%s", img.hash[0:2], img.hash, filepath.Ext(img.originalFilename))
}

func (img *imageRecord) save() error {
	err1 := img.saveToDb()
	err2 := img.saveToDisk()

	if err1 != nil || err2 != nil {
		return fmt.Errorf("Could not save to disk (%v) or to the database(%v),", err2, err1)
	}
	return nil
}

func (img *imageRecord) saveToDb() error {
	_, err := imageInsertStmt.Exec(img.title, img.originalFilename, img.path(), img.height, img.width, img.hash, img.created)
	return err
}

func (img *imageRecord) saveToDisk() error {
	fullPath := filepath.Join(config.UploadPath, img.path())
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(fullPath, img.content, 0644)
}
