package main

import (
	"GoOscilloscopeMusic/wav"
	"os"
)

func main() {
	w := wav.New(2, 8, 480000)
	w.GenerateTone(25, 1, 0.1, 2)
	w.DrawLine(0, 1, 0.1, 1)

	buf := w.Encode()

	file, _ := os.Create("test.wav")
	file.Write(buf)
	defer file.Close()
}