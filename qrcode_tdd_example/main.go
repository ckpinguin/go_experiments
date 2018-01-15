package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"
)

type Version int8 // Top-level type. Alias for int8. Usage example: Version(24)

func GenerateQRCode(w io.Writer, code string, ver Version) error {
	size := ver.PatternSize()
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	return png.Encode(w, img)
}

func (v Version) PatternSize() int {
	return 4*int(v) + 17
}

func main() {
	fmt.Println("Hello QR Code generator")

	file, err := os.Create("qrcode.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	err = GenerateQRCode(file, "0792442222", 1)
	if err != nil {
		log.Fatalln(err)
	}
}
