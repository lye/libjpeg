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

	fout, er := os.Create("test-copy.jpg")
	if er != nil {
		t.Fatal(er)
	}
	defer fout.Close()

	if er := Encode(fout, img, nil); er != nil {
		t.Fatal(er)
	}
}

func TestDecodeFourComponent(t *testing.T) {
	f, er := os.Open("test4c.jpg")
	if er != nil {
		t.Fatal(er)
	}
	defer f.Close()

	_, _, er = image.Decode(f)
	if er != nil {
		t.Fatal(er)
	}
}

func TestDecodeOneComponent(t *testing.T) {
	f, er := os.Open("test1c.jpg")
	if er != nil {
		t.Fatal(er)
	}
	defer f.Close()

	_, _, er = image.Decode(f)
	if er != nil {
		t.Fatal(er)
	}
}
