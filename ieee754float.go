package aiff

// via http://muratnkonar.com/aiff/ , "Converting Extended data to a unsigned long"

type IEEE754Float [10]byte

func (f IEEE754Float) LongValue() uint32 {
	return 0
}
