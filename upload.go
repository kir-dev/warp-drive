package main

import (
	"bytes"
	"crypto"
	_ "crypto/sha1"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"time"
)

func saveFile(title string, file multipart.File, header *multipart.FileHeader) error {
	if title == "" {
		return errors.New("title missing")
	}
	if !checkMimeType(header.Header["Content-Type"]) {
		return errors.New("invalid content-type")
	}

	buff := new(bytes.Buffer)
	if _, err := io.Copy(buff, file); err != nil {
		return err
	}

	imgData := buff.Bytes()
	imgHash, err := hashImage(imgData)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(buff)
	if err != nil {
		return err
	}

	newImg := imageRecord{
		title:            title,
		originalFilename: header.Filename,
		width:            img.Bounds().Size().X,
		height:           img.Bounds().Size().Y,
		hash:             imgHash,
		created:          time.Now(),
		content:          imgData,
	}

	return newImg.save()
}

func checkMimeType(contentType []string) bool {
	allowed := []string{"image/jpeg", "image/gif", "image/png"}

	for _, mime := range allowed {
		for _, ct := range contentType {
			if mime == ct {
				return true
			}
		}
	}

	return false
}

func hashImage(file []byte) (string, error) {
	h := crypto.SHA1.New()
	_, err := h.Write(file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
