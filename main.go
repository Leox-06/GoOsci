package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

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

func (w *wav) makeHeaders(data []byte, SampleRate uint32, BitsPerSample uint16) {
	w.ChunkID = uint32(0x52494646) // "RIFF"
	w.ChunkSize = uint32(36 + len(data))
	w.Format = uint32(0x57415645) // "WAVE"

	w.Subchunk1ID = uint32(0x666d7420) // "fmt "
	w.Subchunk1Size = uint32(16)
	w.AudioFormat = uint16(1)
	w.NumChannels = uint16(2)
	w.SampleRate = SampleRate
	w.ByteRate = (SampleRate * uint32(BitsPerSample) * uint32(w.NumChannels)) / 8
	w.BlockAlign = (w.NumChannels * BitsPerSample) / 8
	w.BitsPerSample = BitsPerSample

	w.Subchunk2ID = uint32(0x64617461) // "data"
	w.Subchunk2Size = uint32(len(data))
	w.data = data
}

func (w wav) createWav(name string) error {
	file, err := os.Create(name + ".wav")
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, w.ChunkID)
	binary.Write(buf, binary.LittleEndian, w.ChunkSize)
	binary.Write(buf, binary.BigEndian, w.Format)

	binary.Write(buf, binary.BigEndian, w.Subchunk1ID)
	binary.Write(buf, binary.LittleEndian, w.Subchunk1Size)
	binary.Write(buf, binary.LittleEndian, w.AudioFormat)
	binary.Write(buf, binary.LittleEndian, w.NumChannels)
	binary.Write(buf, binary.LittleEndian, w.SampleRate)
	binary.Write(buf, binary.LittleEndian, w.ByteRate)
	binary.Write(buf, binary.LittleEndian, w.BlockAlign)
	binary.Write(buf, binary.LittleEndian, w.BitsPerSample)

	binary.Write(buf, binary.BigEndian, w.Subchunk2ID)
	binary.Write(buf, binary.LittleEndian, w.Subchunk2Size)

	binary.Write(buf, binary.LittleEndian, w.data)

	fmt.Printf("%x", buf.Bytes())

	file.Write(buf.Bytes())
	file.Close()

	return nil
}

func main() {
	data := []byte{0x00, 0x00, 0x00, 0x00, 0x24, 0x17, 0x1e, 0xf3, 0x3c, 0x13, 0x3c, 0x14, 0x16, 0xf9, 0x18, 0xf9, 0x34, 0xe7, 0x23, 0xa6, 0x3c, 0xf2, 0x24, 0xf2, 0x11, 0xce, 0x1a, 0x0d}
	var w wav
	w.makeHeaders(data, uint32(48000), uint16(16))
	w.createWav("test")

}
