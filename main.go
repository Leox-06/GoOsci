package main

import (
	"GoOscilloscopeMusic/wav"
	"os"
)

func main() {
	w := wav.New(2, 48000, 8)
	w.Data = []byte{0, 0, 255, 0, 127, 255}
	w.ChangeSpeed(5000)
	buf := w.Encode()

	file, _ := os.Create("test.wav")
	file.Write(buf)
	defer file.Close()
}
