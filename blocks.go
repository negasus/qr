package qr

func mergeBlocks(src [][]byte) []byte {
	var result []byte

	var idx, empty int
	for {
		if empty >= len(src) {
			break
		}
		if idx >= len(src) {
			idx = 0
		}
		if len(src[idx]) == 0 {
			empty++
			idx++
			continue
		}
		empty = 0

		result = append(result, src[idx][0])
		src[idx] = src[idx][1:]

		idx++
	}

	return result
}

func toBlocks(data []byte, level ErrorCorrectionLevel, version int) [][]byte {
	var blocks [][]byte

	count := blocksValues[level][version-1]

	if count == 1 {
		return [][]byte{data}
	}

	a := len(data) / count
	b := len(data) % count

	var idx int

	for i := 0; i < count-b; i++ {
		blocks = append(blocks, data[idx:idx+a])
		idx += a
	}

	for i := 0; i < b; i++ {
		blocks = append(blocks, data[idx:idx+a+1])
		idx += a + 1
	}

	return blocks
}
