// +build gofuzz

package psp

import "bytes"

func Fuzz(data []byte) int {
	cfg, err := DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return 0
	}

	img, err := Decode(bytes.NewReader(data))
	if err != nil {
		if img != nil {
			panic("img != nil on error")
		}
		return 0
	}
	return 1
}
