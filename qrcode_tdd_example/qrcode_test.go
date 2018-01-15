package main

import (
	"bytes"
	"errors"
	"image/png"
	"testing"
)

type ErrorWriter struct{}

func (e *ErrorWriter) Write(b []byte) (int, error) {
	return 0, errors.New("Expected error")
}

func TestGenerateQRCodePropagatesErrors(t *testing.T) {
	w := &ErrorWriter{}
	err := GenerateQRCode(w, "00000000", 1)

	if err == nil || err.Error() != "Expected error" {
		t.Errorf("Error propagation not correct, got %v", err)
	}
}
func TestGenerateQRCodeGeneratesPNG(t *testing.T) {
	buf := new(bytes.Buffer)
	GenerateQRCode(buf, "0792442222", 1)

	if buf.Len() == 0 {
		t.Errorf("No QRCode generated")
	}

	_, err := png.Decode(buf)
	if err != nil {
		t.Errorf("Generated QRCode is not a valid PNG: %s", err)
	}
}

func TestVersionDetermineSize(t *testing.T) {
	table := []struct {
		version  int
		expected int
	}{
		{1, 21},
		{2, 25},
		{6, 41},
		{7, 45},
		{14, 73},
		{40, 177},
	}
	for _, test := range table {
		size := Version(test.version).PatternSize()
		if size != test.expected {
			t.Errorf("Version %2d, expected width %3d but got %3d", test.version, test.expected, size)
		}
	}
}
