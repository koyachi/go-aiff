package aiff

// via http://muratnkonar.com/aiff/ , "Converting Extended data to a unsigned long"

type IEEE754Float80bit [10]byte

func (f *IEEE754Float80bit) LongValue() uint32 {
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

func (f *IEEE754Float80bit) SetLongValue(v uint32) {
	var i uint32
	for i = 0; i < 10; i++ {
		f[i] = 0x00
	}
	exp := v
	exp >>= 1
	for i = 0; i < 32; i++ {
		exp >>= 1
		if exp == 0 {
			break
		}
	}
	f[0] = 0x40 // TODO
	f[1] = byte(i)
	for i = 32; i > 0; i-- {
		if v&0x80000000 == 0x80000000 {
			break
		}
		v <<= 1
	}
	f[2] = uint8((v & 0xFF000000) >> 24)
	f[3] = uint8((v & 0x00FF0000) >> 16)
	f[4] = uint8((v & 0x0000FF00) >> 8)
	f[5] = uint8(v & 0x000000FF)
}
