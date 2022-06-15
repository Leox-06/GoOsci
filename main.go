package main

import (
	"GoOscilloscopeMusic/wav"
	"os"
)

func main() {
	w := wav.New(2, 48000, 8)
	// w.SamplesToData([]float64{0, 0, 1, 0, 1, 1, 0, 1})
	// w.ChangeSpeed(5000)
	w.Transition(0.5, 1, 0.001, 0, 1, 1)
	buf := w.Encode()

	file, _ := os.Create("test.wav")
	file.Write(buf)
	defer file.Close()
}