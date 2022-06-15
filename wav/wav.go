package wav

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

type Wav struct {
	// "RIFF" chunk descriptor
	chunkID   uint32 // "RIFF"
	chunkSize uint32 // 36 + SubChunk2Size
	format    uint32 // "WAVE"

	// "fmt" sub-chunk
	subchunk1ID   uint32 // "fmt "
	subchunk1Size uint32 // sum of the rest subchunk size (2+2+4+4+2+2=16)
	audioFormat   uint16 // 1 for PCM
	NumChannels   uint16 // 2 for stereo
	SampleRate    uint32 // Sample per second
	byteRate      uint32 // ByteRates=(Sample Rate x Bits Per Sample x Channel Numbers)/8
	blockAlign    uint16 // Data block size
	BitsPerSample uint16 // 16bits

	// "data" sub-chunk
	subchunk2ID   uint32 // "data"
	subchunk2Size uint32 // Number of bytes in the data (Sample numbers x Channel numbers x Bits per sample)/8

	Data []byte
}

func New(NumChannels uint16, SampleRate uint32, BitsPerSample uint16) Wav {
	var w Wav

	// "RIFF" chunk descriptor
	w.chunkID = uint32(0x52494646) // "RIFF"
	// chunkSize
	w.format = uint32(0x57415645) // "WAVE"

	// "fmt" sub-chunk
	w.subchunk1ID = uint32(0x666d7420) // "fmt "
	w.subchunk1Size = uint32(16)
	w.audioFormat = uint16(1)
	w.NumChannels = NumChannels
	w.SampleRate = SampleRate
	w.byteRate = (SampleRate * uint32(BitsPerSample) * uint32(NumChannels)) / 8
	w.blockAlign = (NumChannels * BitsPerSample) / 8
	w.BitsPerSample = BitsPerSample

	// "data" sub-chunk
	w.subchunk2ID = uint32(0x64617461) // "data"
	// subchunk2Size

	return w
}

func (w *Wav) Encode() []byte {
	buf := new(bytes.Buffer)

	// headers
	// "RIFF" chunk descriptor
	binary.Write(buf, binary.BigEndian, w.chunkID)
	w.chunkSize = uint32(36 + len(w.Data))
	binary.Write(buf, binary.LittleEndian, w.chunkSize)
	binary.Write(buf, binary.BigEndian, w.format)

	// "fmt" sub-chunk
	binary.Write(buf, binary.BigEndian, w.subchunk1ID)
	binary.Write(buf, binary.LittleEndian, w.subchunk1Size)
	binary.Write(buf, binary.LittleEndian, w.audioFormat)
	binary.Write(buf, binary.LittleEndian, w.NumChannels)
	binary.Write(buf, binary.LittleEndian, w.SampleRate)
	binary.Write(buf, binary.LittleEndian, w.byteRate)
	binary.Write(buf, binary.LittleEndian, w.blockAlign)
	binary.Write(buf, binary.LittleEndian, w.BitsPerSample)

	// "data" sub-chunk
	binary.Write(buf, binary.BigEndian, w.subchunk2ID)
	w.subchunk2Size = uint32(len(w.Data))
	binary.Write(buf, binary.LittleEndian, w.subchunk2Size)

	// data
	binary.Write(buf, binary.LittleEndian, w.Data)

	return buf.Bytes()
}

func (w *Wav) ChangeSpeed(multiplier uint) {
	var newData []byte
	for i := 0; i < len(w.Data); i += int(w.NumChannels) {
		for m := 0; m < int(multiplier); m++ {
			newData = append(newData, w.Data[i], w.Data[i+1])
		}
	}
	w.Data = newData
}

func (w *Wav) SamplesToData(samples []float64, start float64, channel ...int) {
	for _, v := range channel {
		if v > int(w.NumChannels) {
			panic("channel doesn't exist")
		}
	}

	w.Data = make([]byte, len(samples)*int(w.NumChannels)*8/int(w.BitsPerSample))
	for i := 0; i < len(samples); i++ {
		sampleBits := byte(samples[i] * (math.Pow(2, float64(w.BitsPerSample)) - 1))
		for c := 0; c < int(w.NumChannels); c++ {
			if len(channel) > c {
				w.Data[i*int(w.NumChannels)+channel[c]-1] = sampleBits
			}
		}
	}
}

func (w *Wav) GenerateTone(frequency, amplitude, start, duration float64, channel ...int) {
	if amplitude < 0 || amplitude > 1 {
		panic("the amplitude must be between 0 and 1")
	}

	var samples []float64
	for i := 0.0; i < duration; i += 1 / float64(w.SampleRate) {
		sample := (amplitude*math.Sin(i*2*math.Pi*frequency) + 1) / 2
		samples = append(samples, sample)
	}

	w.SamplesToData(samples, start, channel...)
}

func (w *Wav) Transition(a, b, s, start, duration float64, channel ...int) {
	samples := []float64{a}
	for i := 0.0; i < duration; i += 1 / float64(w.SampleRate) {
		samples = append(samples, lerp(samples[len(samples)-1], b, s))
	}
	fmt.Println(samples)

	w.SamplesToData(samples, start, channel...)
}

func lerp(a, b, s float64) float64 {
	r := float64(b - a)
	return a + (r * s)
}
