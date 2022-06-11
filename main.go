package main

import "os"

type wav struct {
	// "RIFF" chunk descriptor
	ChunkID   uint32 // "RIFF"
	ChunkSize uint32 // 36 + SubChunk2Size
	Format    uint32 // "WAVE"

	// "fmt" sub-chunk
	Subchunk1ID   uint32 // "fmt "
	Subchunk1Size uint32 // sum of the rest subchunk size (2+2+4+4+2+2=16)
	AudioFormat   uint16 // 1 for PCM
	NumChannels   uint16 // 2 for stereo
	SampleRate    uint32 // 48000 Hz (80 BB 00 00)
	ByteRate      uint32 // ByteRates=(Sample Rate x Bits Per Sample x Channel Numbers)/8
	BlockAlign    uint16 // Data block size
	BitsPerSample uint16 // 16bits

	// "data" sub-chunk
	Subchunk2ID   uint32 // "data"
	Subchunk2Size uint32 // Number of bytes in the data (Sample numbers x Channel numbers x Bits per sample)/8

	data []byte
}

func (w *wav) makeHeaders(SampleRate uint32, BitsPerSample uint16, data []byte) {
	if SampleRate == 0 {
		SampleRate = uint32(48000)
	}
	if BitsPerSample == 0 {
		BitsPerSample = uint16(16)
	}

	w.ChunkID = uint32(0x52494646) // "RIFF"
	w.ChunkSize = 36 + w.Subchunk2Size
	w.Format = uint32(0x57415645) // "WAVE"

	w.Subchunk1ID = uint32(0x666d7420)
	w.Subchunk1Size = uint32(16)
	w.AudioFormat = uint16(1)
	w.NumChannels = uint16(2)
	w.SampleRate = SampleRate
	w.ByteRate = (SampleRate * uint32(BitsPerSample) * uint32(w.NumChannels)) / 8
	w.BlockAlign = uint16(len(w.data))
	w.BitsPerSample = BitsPerSample

	w.Subchunk1ID = uint32(0x64617461)
	w.Subchunk2Size = uint32(len(data))
	w.data = data
}

func main() {

	data := []byte{255, 0, 255, 0, 255, 0, 255, 0}

	file, _ := os.Create("file.wav")
	file.Write(data)
	file.Close()
}
