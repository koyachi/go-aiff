package aiff

import (
	"io"
)

type AIFFFormat struct {
	//	FormType        uint16
	NumChannels     uint16
	NumSampleFrames uint32
	SampleSize      uint16
	SampleRate      [10]byte
	//SampleRate IEEE754Float
}

type AIFFData struct {
	io.Reader
	Size uint32
	//offset uint32
	//blockSize uint32
	pos uint32
}

type Sample struct {
	Offset    uint16
	BlockSize uint16
	Values    [2]int
}
