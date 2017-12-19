package main

import (
	"bytes"
	"image/png"
	"testing"
)

func TestGenerateQRCodeGeneratesPNG(t *testing.T) {
	buf := new(bytes.Buffer)
	GenerateQRCode(buf, "0792442222")

	if buf.Len() == 0 {
		t.Errorf("No QRCode generated")
	}

	_, err := png.Decode(buf)
	if err != nil {
		t.Errorf("Generated QRCode is not a valid PNG: %s", err)
	}
}
