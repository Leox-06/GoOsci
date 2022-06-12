package main

import (
	"GoOscilloscopeMusic/wav"
	"os"
)

func main() {
	w := wav.New(1, 4000, 8)
	w.Data = append(w.Data, []byte{127, 127, 127, 255, 255, 255,127, 127, 127, 0, 0, 0}...)
	buf := w.Encode()

	file, _ := os.Create("test.wav")
	file.Write(buf)
	defer file.Close()
}
