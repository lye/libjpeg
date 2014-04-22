package libjpeg

import (
	"image"
	"os"
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

	fout, er := os.Create("test2.jpg")
	if er != nil {
		t.Fatal(er)
	}
	defer fout.Close()

	if er := Encode(fout, img, nil); er != nil {
		t.Fatal(er)
	}
}
