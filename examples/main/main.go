package main

import "github.com/negasus/qr"

func main() {
	err := qr.SaveImage([]byte("https://negasus.dev"), "qr.png")
	if err != nil {
		panic(err)
	}
}
