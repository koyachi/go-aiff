package aiff

import (
	"bufio"
	"encoding/binary"
	"errors"
	iff "github.com/koyachi/go-iff"
	"math"
)

type Reader struct {
	r        *iff.Reader
	iffChunk *iff.IFFChunk
	format   *AIFFFormat
	*AIFFData
	lastChunkIndex map[string]int
}

func NewReader(r iff.IFFReader) *Reader {
	iffReader := iff.NewReader(r)
	return &Reader{r: iffReader}
}

func (r *Reader) Format() (format *AIFFFormat, err error) {
	if r.format == nil {
		format, err = r.readFormat()
		if err != nil {
			return
		}
		r.format = format
	} else {
		format = r.format
	}

	return
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if r.AIFFData == nil {
		data, err := r.readData()
		if err != nil {
			return n, err
		}
		r.AIFFData = data
	}

	return r.AIFFData.Read(p)
}

func (r *Reader) ReadSamples(params ...uint32) (samples []Sample, err error) {
	var bytes []byte
	var numSamples, n, buffered int

	if len(params) > 0 {
		numSamples = int(params[0])
	} else {
		numSamples = 2048
	}

	format, err := r.Format()
	if err != nil {
		return
	}

	numChannels := int(format.NumChannels)
	sampleSize := int(format.SampleSize)
	var dataOffset uint32
	var blockSize uint32

	dataOffset, err = bytesToBEUInt32(r)
	if err != nil {
		return
	}

	blockSize, err = bytesToBEUInt32(r)
	if err != nil {
		return
	}
	if blockSize != 0 {
		panic("TODO: implement block alignemtn")
	}

	buffered = 0
	n = int(dataOffset)
	for buffered < n {
		bytes = make([]byte, dataOffset)
		n, err = r.Read(bytes)
		if err != nil {
			return
		}
		r.AIFFData.pos += uint32(n)
		buffered += n
	}

	samples = make([]Sample, numSamples)

	buffered = 0
	n = numSamples * sampleSize * numChannels
	for buffered < n {
		// TODO:FIX
		bytes = make([]byte, numSamples*sampleSize*numChannels)
		n, err = r.AIFFData.Read(bytes)
		if err != nil {
			return
		}
		r.AIFFData.pos += uint32(numSamples) * uint32(sampleSize) * uint32(numChannels)

		offset := 0
		for i := 0; i < numSamples; i++ {
			if sampleSize == 16 {
				for j := 0; j < int(numChannels); j++ {
					soffset := offset + (j * numChannels)
					samples[i].Values[j] = int((int16(bytes[soffset]) << 8) + int16(bytes[soffset+1]))
				}
			} else {
				for j := 0; j < int(numChannels); j++ {
					samples[i].Values[j] = int(bytes[offset+j])
				}
			}

			offset += sampleSize / 8 * numChannels
		}
		buffered += n
	}

	return
}

func (r *Reader) IntValue(sample Sample, channel uint) int {
	return sample.Values[channel]
}

func (r *Reader) FloatValue(sample Sample, channel uint) float64 {
	// XXX
	return float64(r.IntValue(sample, channel)) / math.Pow(2, float64(r.format.SampleSize))
}

func (r *Reader) readFormat() (fmt *AIFFFormat, err error) {
	var iffChunk *iff.IFFChunk

	fmt = new(AIFFFormat)

	if r.iffChunk == nil {
		// TODO: "RIFF"で比較してるのでエラーになるのをどうにかする(go-riff的には正しい) go-iff?
		iffChunk, err = r.r.Read()
		if err != nil {
			return
		}

		r.iffChunk = iffChunk
	} else {
		iffChunk = r.iffChunk
	}

	commonChunk := findChunk(iffChunk, "COMM")
	if commonChunk == nil {
		err = errors.New("Common chunk is not found")
		return
	}

	err = binary.Read(commonChunk, binary.BigEndian, fmt)
	if err != nil {
		return
	}

	return
}

func (r *Reader) readData() (data *AIFFData, err error) {
	var iffChunk *iff.IFFChunk

	if r.iffChunk == nil {
		iffChunk, err = r.r.Read()
		if err != nil {
			return
		}

		r.iffChunk = iffChunk
	} else {
		iffChunk = r.iffChunk
	}

	soundDataChunk := findChunk(iffChunk, "SSND")
	if soundDataChunk == nil {
		err = errors.New("Sound Data chunk is not found")
		return
	}

	data = &AIFFData{bufio.NewReader(soundDataChunk), soundDataChunk.ChunkSize, 0}

	return
}

func findChunk(iffChunk *iff.IFFChunk, id string) (chunk *iff.Chunk) {
	for _, ch := range iffChunk.Chunks {
		if string(ch.ChunkID[:]) == id {
			chunk = ch
			break
		}
	}

	return
}

/*
func (r *Reader) findChunkNext(id string) (chunk *iff.Chunk) {
	for i, ch := range r.iffChunk.Chunks {
		if string(ch.ChunkID[:]) == id {
			chunk = ch
			break
		}
	}

	if r.lastChunkIndex == nil {
		r.lastChunkIndex = make(map[string]int)
	}
	r.lastChunkIndex[id] = i

	return
}
*/

func bytesToBEUInt32(r *Reader) (uint32, error) {
	tmp := make([]byte, 4)
	n, err := r.Read(tmp)
	if err != nil {
		return 0, err
	}
	if n != 4 {
		return 0, errors.New("can't read dataOffset")
	}
	r.AIFFData.pos += uint32(n)
	return uint32(tmp[0])<<24 +
		uint32(tmp[1])<<16 +
		uint32(tmp[2])<<8 +
		uint32(tmp[3]), nil
}
