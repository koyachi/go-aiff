package aiff

// via http://muratnkonar.com/aiff/ , "Converting Extended data to a unsigned long"

type IEEE754Float80bit [10]byte

func (f IEEE754Float80bit) LongValue() uint32 {
	mantissa := uint32(f[2])<<24 | uint32(f[3])<<16 | uint32(f[4])<<8 | uint32(f[5])
	exp := 30 - f[1]
	var last uint32
	for ; exp > 0; exp-- {
		last = mantissa
		mantissa >>= 1
	}
	if last&0x00000001 == 0x00000001 {
		mantissa++
	}
	return mantissa
}
