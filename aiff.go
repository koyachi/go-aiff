package aiff

import (
	"io"
)

type AIFFFormat struct {
	//	FormType        uint16
	NumChannels     uint16
	NumSampleFrames uint32
	SampleSize      uint16
	SampleRate      IEEE754Float80bit
}

type AIFFData struct {
	io.Reader
	Size uint32
	//offset uint32
	//blockSize uint32
	pos uint32
}

type Sample struct {
	Values [2]int
}
