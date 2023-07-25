package qr

type bitwise struct {
	data []byte
}

func (b *bitwise) padZero() {
	l := len(b.data)
	if l%8 == 0 {
		return
	}
	for i := 0; i < 8-l%8; i++ {
		b.data = append(b.data, 0)
	}
}

func (b *bitwise) addBytes(data []byte) {
	for _, sym := range data {
		var arr []byte
		for i := 0; i < 8; i++ {
			arr = append(arr, (sym>>uint(i))&1)
		}
		for j := len(arr) - 1; j >= 0; j-- {
			b.data = append(b.data, arr[j])
		}
	}
}

func (b *bitwise) addNumber(n int, bits int) {
	var v []byte
	for j := 0; j < bits; j++ {
		v = append(v, byte((n>>uint(j))&1))
	}
	for j := len(v) - 1; j >= 0; j-- {
		b.data = append(b.data, v[j])
	}
}

func (b *bitwise) addNumberToStart(n int, bits int) {
	for j := 0; j < bits; j++ {
		b.data = append([]byte{byte((n >> uint(j)) & 1)}, b.data...)
	}
}

func (b *bitwise) addPaddingToBytesLen(toLen int) {
	paddingNumbers := []byte{0b11101100, 0b00010001}

	var idx int
	for {
		if len(b.data)/8 >= toLen {
			return
		}
		e := paddingNumbers[idx]
		idx++
		if idx >= len(paddingNumbers) {
			idx = 0
		}
		b.addNumber(int(e), 8)
	}
}

func (b *bitwise) bytes() []byte {
	var data []byte
	for i := 0; i < len(b.data); i += 8 {
		var e byte
		for j := 0; j < 8; j++ {
			e += b.data[i+j] << (7 - j)
		}
		data = append(data, e)
	}
	return data
}

func (b *bitwise) len() int {
	return len(b.data)
}
