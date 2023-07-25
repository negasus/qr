package qr

type walker struct {
	x      int
	y      int
	size   int
	data   []byte
	dy     int
	shiftX bool
}

func newWalker(size int, data []byte) *walker {
	c := &walker{
		x:      size - 2,
		y:      size,
		size:   size,
		data:   make([]byte, len(data)),
		dy:     -1,
		shiftX: true,
	}
	copy(c.data, data)
	return c
}

func (w *walker) next() (int, int) {
	for {
		if w.x < 0 || w.y < 0 {
			return -1, -1
		}

		w.shiftX = !w.shiftX

		if w.shiftX {
			w.x--
			if w.data[w.y*w.size+w.x] != 0xFF {
				continue
			}
			return w.x, w.y
		}

		// skip vertical timing pattern
		if w.x == 5 && w.y == 9 {
			w.shiftX = false
			if w.data[w.y*w.size+w.x] != 0xFF {
				continue
			}
			return w.x, w.y
		}

		w.x++
		w.y += w.dy

		if w.y < 0 {
			w.y = 0
			w.dy = -w.dy
			w.x -= 2
		}
		if w.y >= w.size {
			w.y = w.size - 1
			w.dy = -w.dy
			w.x -= 2
		}

		if w.data[w.y*w.size+w.x] != 0xFF {
			continue
		}

		return w.x, w.y
	}
}
