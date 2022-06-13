package main

import (
	"GoOscilloscopeMusic/wav"
	"os"
)

func main() {
	w := wav.New(1, 4000, 8)
	w.GenerateTone(5, 1, 1)
	buf := w.Encode()

	file, _ := os.Create("test.wav")
	file.Write(buf)
	defer file.Close()
}
