package libjpeg

import (
	"os"
	"image"
	"image/png"
	"testing"
)

func TestReencode(t *testing.T) {
	fin, er := os.Open("test.jpg")
	if er != nil {
		t.Fatal(er)
	}
	defer fin.Close()

	img, _, er := image.Decode(fin)
	if er != nil {
		t.Fatal(er)
	}

	fout, er := os.Create("test.png")
	if er != nil {
		t.Fatal(er)
	}
	defer fout.Close()

	if er := png.Encode(fout, img) ; er != nil {
		t.Fatal(er)
	}
}
