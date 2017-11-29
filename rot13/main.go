package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (rot *rot13Reader) Read(p []byte) (n int, err error) {
	n, err = rot.r.Read(p)
	for i := 0; i < n; i++ {
		var o byte
		switch {
		case p[i] >= byte('A') && p[i] <= byte('M'):
			{
				o = p[i] + 13
			}
		case p[i] >= byte('a') && p[i] <= byte('m'):
			{
				o = p[i] + 13
			}
		case p[i] >= byte('N') && p[i] <= byte('Z'):
			{
				o = p[i] - 13
			}
		case p[i] >= byte('N') && p[i] <= byte('z'):
			{
				o = p[i] - 13
			}
		default:
			{
				o = p[i]
			}
		}
		p[i] = o
	}

	return
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
