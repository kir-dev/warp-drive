package main

import (
	"testing"
)

func TestCannotResizeGifImages(t *testing.T) {
	_, err := resizeImage("grumpy-cat.gif", 500)

	if err == nil || err.Error() != "resizing gifs are not supported" {
		t.Error("resizing gifs should throw an error")
	}
}
