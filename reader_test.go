package aiff

import (
	"io/ioutil"
	"math"
	"testing"
)

func TestRead(t *testing.T) {
	//	blockAlign := 4

	file, err := fixtureFile("a.aiff")

	if err != nil {
		t.Fatalf("Failed to open fixture file")
	}

	reader := NewReader(file)
	fmt, err := reader.Format()
	if err != nil {
		t.Fatal(err)
	}

	//	if fmt.AudioFormat != AudioFormatPCM {
	//		t.Fatalf("Audio format is invalid: %d", fmt.AudioFormat)
	//	}

	if fmt.NumChannels != 2 {
		t.Fatalf("Number of channels is invalid: %d", fmt.NumChannels)
	}

	// TODO
	//	if fmt.SampleRate != 44100 {
	//		t.Fatalf("Sample rate is invalid: %d", fmt.SampleRate)
	//	}

	//	if fmt.SampleRate != 44100*4 {
	//		t.Fatalf("Byte rate is invalid: %d", fmt.ByteRate)
	//	}

	//	if int(fmt.BlockAlign) != blockAlign {
	//		t.Fatalf("Block align is invalid: %d", fmt.BlockAlign)
	//	}

	if fmt.SampleSize != 16 {
		t.Fatalf("Sample size is invalid: %d", fmt.SampleSize)
	}

	numSamples := 1
	samples, err := reader.ReadSamples(uint32(numSamples))

	if len(samples) != numSamples {
		t.Fatalf("Length of samples is invalid: %d", len(samples))
	}

	sample := samples[0]

	if reader.IntValue(sample, 0) != 318 {
		t.Fatalf("Value is invalid: %d", reader.IntValue(sample, 0))
	}

	if reader.IntValue(sample, 1) != 289 {
		t.Fatalf("Value is invalid: %d", reader.IntValue(sample, 1))
	}

	if math.Abs(reader.FloatValue(sample, 0)-0.004852) > 0.0001 {
		t.Fatalf("Value is invalid: %f", reader.FloatValue(sample, 0))
	}

	if math.Abs(reader.FloatValue(sample, 1)-0.004409) > 0.0001 {
		t.Fatalf("Value is invalid: %d", reader.FloatValue(sample, 1))
	}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	chunkHeaderPlusOneSampleSize := 4 + 4 + (fmt.NumChannels * fmt.SampleSize * uint16(numSamples))
	if len(bytes) != int(reader.AIFFData.Size)-int(chunkHeaderPlusOneSampleSize) {
		t.Fatalf("Data size is invalid: %d, %d", len(bytes), reader.AIFFData.Size)
	}

	t.Logf("Data size: %d", len(bytes))
}
