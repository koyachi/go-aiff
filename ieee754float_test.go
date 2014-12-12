package aiff

import (
	"bytes"
	"testing"
)

func TestIEEE754Float80bit(t *testing.T) {
	var data [10]byte
	data = [10]byte{0x40, 0x0e, 0xac, 0x44, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	var f IEEE754Float80bit
	f = data
	if f.LongValue() != 44100 {
		t.Fatalf("LongValue() is invalid: %d", f.LongValue())
	}

	var data2 [10]byte
	data2 = [10]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	var f2 IEEE754Float80bit
	f2 = data2
	f2.SetLongValue(44100)
	expected := [10]byte{0x40, 0x0e, 0xac, 0x44, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	if bytes.Compare(f2[:], expected[:]) != 0 {
		t.Fatalf("SetLongValue(44100) is invalid: %v", f2)
	}
}
