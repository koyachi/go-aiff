package aiff

import (
	"encoding/binary"
	iff "github.com/koyachi/go-iff"
	"io"
)

type Writer struct {
	io.Writer
	Format *AIFFFormat
}

func NewWriter(w io.Writer, numSamples uint32, numChannels uint16, sampleRate uint32, sampleSize uint16) (writer *Writer) {
	numSampleFrames := numSamples // TODO
	sampleRate80bit := IEEE754Float80bit{}
	sampleRate80bit.SetLongValue(sampleRate)
	format := &AIFFFormat{numChannels, uint32(numSampleFrames), sampleSize, sampleRate80bit}

	dataSize := numSamples * uint32(numChannels) * uint32(sampleSize) // TODO
	iffSize := 4 + 8 + 18 + 8 + dataSize                              // TODO
	iffWriter := iff.NewWriter(w, []byte("AIFF"), iffSize)

	writer = &Writer{iffWriter, format}
	iffWriter.WriteChunk([]byte("COMM"), 18, func(w io.Writer) {
		binary.Write(w, binary.BigEndian, format)
	})
	iffWriter.WriteChunk([]byte("SSND"), dataSize, func(w io.Writer) {
		var dataOffset uint32 = 0
		var blockSize uint32 = 0
		binary.Write(w, binary.BigEndian, dataOffset)
		binary.Write(w, binary.BigEndian, blockSize)
	})

	return writer
}

func (w *Writer) WriteSamples(samples []Sample) (err error) {
	sampleSize := w.Format.SampleSize
	numChannels := w.Format.NumChannels

	var i uint16
	for _, sample := range samples {
		if sampleSize == 16 {
			for i = 0; i < numChannels; i++ {
				err = binary.Write(w, binary.BigEndian, int16(sample.Values[i]))
				if err != nil {
					return
				}
			}
		} else {
			for i = 0; i < numChannels; i++ {
				err = binary.Write(w, binary.BigEndian, int8(sample.Values[i]))
				if err != nil {
					return
				}
			}
		}
	}

	return
}
