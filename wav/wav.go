package wav

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type Wav struct {
	NumChannels   int // number of cannels
	SampleRate    int // Sample per second
	BitsPerSample int // bits per sample

	Data []byte
}

func New(numChannels, sampleRate, bitsPerSample int) *Wav{
	return &Wav{NumChannels: numChannels, SampleRate: sampleRate, BitsPerSample: bitsPerSample}
}

func (w *Wav) Encode() []byte {
	buf := new(bytes.Buffer)

	// headers
	// "RIFF" chunk descriptor
	RIFF := uint32(0x52494646) // "RIFF"
	binary.Write(buf, binary.BigEndian, RIFF)

	chunkSize := uint32(36 + len(w.Data))
	binary.Write(buf, binary.LittleEndian, chunkSize)

	format := uint32(0x57415645) // "WAVE"
	binary.Write(buf, binary.BigEndian, format)

	// "fmt" sub-chunk
	subchunk1ID := uint32(0x666d7420) // "fmt "
	binary.Write(buf, binary.BigEndian, subchunk1ID)

	subchunk1Size := uint32(16)
	binary.Write(buf, binary.LittleEndian, subchunk1Size)

	audioFormat := uint16(1)
	binary.Write(buf, binary.LittleEndian, audioFormat)

	numChannels := uint16(w.NumChannels)
	binary.Write(buf, binary.LittleEndian, numChannels)

	sampleRate := uint32(w.SampleRate)
	binary.Write(buf, binary.LittleEndian, sampleRate)

	bitsPerSample := uint16(w.BitsPerSample)

	byteRate := (sampleRate * uint32(bitsPerSample) * uint32(numChannels)) / 8
	binary.Write(buf, binary.LittleEndian, byteRate)

	blockAlign := (numChannels * bitsPerSample) / 8
	binary.Write(buf, binary.LittleEndian, blockAlign)

	binary.Write(buf, binary.LittleEndian, bitsPerSample)

	// "data" sub-chunk
	subchunk2ID := uint32(0x64617461) // "data"
	binary.Write(buf, binary.BigEndian, subchunk2ID)

	subchunk2Size := uint32(len(w.Data))
	binary.Write(buf, binary.LittleEndian, subchunk2Size)

	// data
	binary.Write(buf, binary.LittleEndian, w.Data)

	return buf.Bytes()
}

func (w *Wav) ChangeSpeed(divider uint) {
	var newData []byte
	for i := 0; i < len(w.Data); i += w.NumChannels {
		for m := 0; m < int(divider); m++ {
			newData = append(newData, w.Data[i], w.Data[i+1])
		}
	}
	w.Data = newData
}

func (w *Wav) SamplesToData(samples []float64, channels ...int) {
	for _, v := range channels {
		if v > w.NumChannels {
			panic("channel doesn't exist")
		}
	}

	if len(samples) > len(w.Data) {
		length := (len(samples) - len(w.Data)) * w.NumChannels * (w.BitsPerSample / 8)
		w.Data = append(w.Data, make([]byte, length)...)
	}

	for i := 0; i < len(samples); i++ {
		sampleBits := byte(samples[i] * (math.Pow(2, float64(w.BitsPerSample)) - 1))
		for c := 0; c < w.NumChannels; c++ {
			if len(channels) > c {
				w.Data[i*w.NumChannels+channels[c]-1] = sampleBits
			}
		}
	}
}

func (w *Wav) GenerateTone(frequency, amplitude, duration float64, channel ...int) {
	if amplitude < 0 || amplitude > 1 {
		panic("the amplitude must be between 0 and 1")
	}

	var samples []float64
	for i := 0.0; i < duration; i += 1 / float64(w.SampleRate) {
		sample := (amplitude*math.Sin(i*2*math.Pi*frequency) + 1) / 2
		samples = append(samples, sample)
	}

	w.SamplesToData(samples, channel...)
}

func (w *Wav) DrawLine(begin, end, duration float32, channels ...int) {
	tween := gween.New(begin, end, duration, ease.Linear)

	samples := []float64{}
	var isFinished bool

	for !isFinished {
		var current float32
		current, isFinished = tween.Update(1 / float32(w.SampleRate))
		samples = append(samples, float64(current))
	}

	w.SamplesToData(samples, channels...)
}
