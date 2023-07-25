package qr

type builder struct {
	size   int
	canvas []byte
}

func newBuilder(size int) *builder {
	r := &builder{
		size:   size,
		canvas: make([]byte, size*size),
	}

	for i := 0; i < len(r.canvas); i++ {
		r.canvas[i] = 0xFF
	}
	return r
}

func (r *builder) build(data []byte, maskID int, level ErrorCorrectionLevel, version int) []byte {
	maskValue := maskValues[level][maskID]

	r.addPositions()
	r.addTiming()

	r.addMaskAndCorrectionLevel(maskValue)
	if version > 1 {
		r.addAlignments(version)
	}
	if version > 6 {
		r.addVersion(version)
	}

	r.addData(data, maskID)

	return r.canvas
}

func (r *builder) addVersion(version int) {
	c := versions[version-7]
	var idx int
	for y := r.size - 9; y >= r.size-11; y-- {
		for x := 5; x >= 0; x-- {
			if c>>idx&1 == 1 {
				r.set(x, y, 1)
				r.set(y, x, 1)
			} else {
				r.set(x, y, 0)
				r.set(y, x, 0)
			}
			idx++
		}
	}
}

func (r *builder) addAlignments(version int) {
	data := alignmentPositions[version-1]

	for x := 0; x < len(data); x++ {
		for y := 0; y < len(data); y++ {
			if version > 6 && ((x == 0 && y == 0) || (x == 0 && y == len(data)-1) || (x == len(data)-1 && y == 0)) {
				continue
			}

			r.set(data[x], data[y], 1)
			r.rect(data[x]-1, data[y]-1, data[x]+1, data[y]+1, 0)
			r.rect(data[x]-2, data[y]-2, data[x]+2, data[y]+2, 1)
		}
	}
}

func (r *builder) set(x, y int, v byte) {
	if x < 0 || y < 0 || x >= r.size || y >= r.size {
		return
	}
	r.canvas[y*r.size+x] = v
}

func (r *builder) vLine(x, y1, y2 int, v byte) {
	for y := y1; y <= y2; y++ {
		r.set(x, y, v)
	}
}

func (r *builder) hLine(x1, x2, y int, v byte) {
	for x := x1; x <= x2; x++ {
		r.set(x, y, v)
	}
}

func (r *builder) rect(x1, y1, x2, y2 int, v byte) {
	r.hLine(x1, x2, y1, v)
	r.hLine(x1, x2, y2, v)
	r.vLine(x1, y1, y2, v)
	r.vLine(x2, y1, y2, v)
}

func (r *builder) drawPosition(x, y int) {
	r.rect(x-1, y-1, x+7, y+7, 0)
	r.rect(x, y, x+6, y+6, 1)
	r.rect(x+1, y+1, x+5, y+5, 0)
	r.rect(x+2, y+2, x+4, y+4, 1)
	r.set(x+3, y+3, 1)
}

func (r *builder) addPositions() {
	r.drawPosition(0, 0)        // top left
	r.drawPosition(r.size-7, 0) // top right
	r.drawPosition(0, r.size-7) // bottom left
}

func (r *builder) addTiming() {
	for i := 8; i < r.size-8; i++ {
		if i%2 == 0 {
			r.set(i, 6, 1)
			r.set(6, i, 1)
		} else {
			r.set(i, 6, 0)
			r.set(6, i, 0)
		}
	}
}

type point struct {
	x, y int
}

func (r *builder) addMaskAndCorrectionLevel(mask int) {
	path1 := []point{
		{0, 8}, {1, 8}, {2, 8}, {3, 8}, {4, 8}, {5, 8}, {7, 8}, {8, 8},
		{8, 7}, {8, 5}, {8, 4}, {8, 3}, {8, 2}, {8, 1}, {8, 0},
	}

	y := r.size - 1
	x := r.size - 1

	path2 := []point{
		{8, y}, {8, y - 1}, {8, y - 2}, {8, y - 3}, {8, y - 4}, {8, y - 5}, {8, y - 6},
		{x - 7, 8}, {x - 6, 8}, {x - 5, 8}, {x - 4, 8}, {x - 3, 8}, {x - 2, 8}, {x - 1, 8}, {x, 8},
	}

	for i := 0; i < 15; i++ {
		c := byte(0)
		if mask>>i&1 == 1 {
			c = 1
		}
		r.set(path1[14-i].x, path1[14-i].y, c)
		r.set(path2[14-i].x, path2[14-i].y, c)
	}

	r.set(8, r.size-8, 1)
}

func (r *builder) addData(data []byte, mask int) {
	c := newWalker(r.size, r.canvas)

	for _, v := range data {
		for j := 0; j < 8; j++ {
			x, y := c.next()

			if x < 0 || y < 0 {
				panic("x<0 || y<0")
			}

			isSet := v<<j&0x80 == 0x80

			if r.mask(mask, x, y) == 0 {
				isSet = !isSet
			}

			if isSet {
				r.set(x, y, 1)
			}
		}
	}
}

func (r *builder) mask(v, x, y int) int {
	switch v {
	case 0:
		return (x + y) % 2
	case 1:
		return y % 2
	case 2:
		return x % 3
	case 3:
		return (x + y) % 3
	case 4:
		return (x/3 + y/2) % 2
	case 5:
		return (x*y)%2 + (x*y)%3
	case 6:
		return ((x*y)%2 + (x*y)%3) % 2
	case 7:
		return ((x*y)%3 + (x+y)%2) % 2
	}

	return -1
}
