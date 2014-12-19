package aiff

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	outfile, err := ioutil.TempFile("/tmp", "outfile")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		outfile.Close()
		os.Remove(outfile.Name())
	}()

	var numSamples uint32 = 2
	var numChannels uint16 = 2
	var sampleRate uint32 = 44100
	var sampleSize uint16 = 16

	writer := NewWriter(outfile, numSamples, numChannels, sampleRate, sampleSize)
	samples := make([]Sample, numSamples)

	samples[0].Values[0] = 32767
	samples[0].Values[1] = -32768
	samples[1].Values[0] = 123
	samples[1].Values[1] = -123

	err = writer.WriteSamples(samples)
	if err != nil {
		t.Fatal(err)
	}

	outfile.Close()
	file, err := os.Open(outfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		file.Close()
		os.Remove(outfile.Name())
	}()

	reader := NewReader(file)
	if err != nil {
		t.Fatal(err)
	}

	fmt, err := reader.Format()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, numChannels, fmt.NumChannels)
	assert.Equal(t, sampleRate, fmt.SampleRate.LongValue())
	assert.Equal(t, sampleSize, fmt.SampleSize)

	samples, err = reader.ReadSamples(2)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, len(samples))
	assert.Equal(t, 32767, samples[0].Values[0])
	assert.Equal(t, -32768, samples[0].Values[1])
	assert.Equal(t, 123, samples[1].Values[0])
	assert.Equal(t, -123, samples[1].Values[1])
}
