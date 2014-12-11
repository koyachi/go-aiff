package aiff

import (
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
}
