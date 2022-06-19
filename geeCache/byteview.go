package geeCache

type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return v.Len()
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

// the purpose is to prevent the cache value from being  modified
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
