package psp

import (
	"fmt"
	"image/png"
	"os"
	"testing"
)

func TestDecode(t *testing.T) {
	// f, err := os.Open("3000939.psp")
	// f, err := os.Open("v13-rgb-lz77.pspimage")
	// f, err := os.Open("v13-rgb16-lz77.pspimage")
	// f, err := os.Open("v13-rgba-lz77.pspimage")
	// f, err := os.Open("v13-bw-lz77.pspimage")
	// f, err := os.Open("v13-gray16-lz77.pspimage")
	// f, err := os.Open("v13-paletted-lz77.pspimage")
	// f, err := os.Open("left-corners.pspimage")
	// f, err := os.Open("v12-paletted-lz77.pspimage")
	// f, err := os.Open("v10b-paletted-lz77.pspimage")
	// f, err := os.Open("v10a-paletted-lz77.pspimage")
	// f, err := os.Open("v9-paletted-lz77.pspimage")
	// f, err := os.Open("v8-paletted-lz77.pspimage")
	f, err := os.Open("../testdata/v7-paletted-lz77.pspimage")
	// f, err := os.Open("../testdata/v6-paletted-lz77.pspimage")
	// f, err := os.Open("../testdata/v5-paletted-lz77.pspimage")
	// f, err := os.Open("v4-paletted-lz77.pspimage")
	// f, err := os.Open("v3-paletted-lz77.pspimage")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	img, err := Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	fo, err := os.Create("test.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(fo, img); err != nil {
		t.Fatal(err)
	}
}

func TestDecodeConfig(t *testing.T) {
	// f, err := os.Open("3000939.psp")
	f, err := os.Open("Nibbler.pspimage")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	config, err := DecodeConfig(f)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", config)
}
