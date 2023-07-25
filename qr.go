package qr

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
)

type Renderer interface {
	Render(data []byte) (image.Image, error)
}

func EncodeBytes(errorCorrectionLevel ErrorCorrectionLevel, maskID int, value []byte, renderer Renderer, destFile string) error {
	b := &bitwise{}
	b.addBytes(value)

	var version, ls int

	version = calcVersion(errorCorrectionLevel, len(value)*8)
	if version < 0 {
		return fmt.Errorf("data is too big")
	}

	for {
		ls = getLengthFieldSize(EncodeTypeByte, version)

		if len(value)*8+ls+4 <= maxBitsValues[errorCorrectionLevel][version-1] {
			break
		}

		version++
		if version > 40 {
			return fmt.Errorf("data is too big")
		}
	}

	b.addNumberToStart(len(value), ls)
	b.addNumberToStart(int(EncodeTypeByte), 4)
	b.padZero()

	maxBits := maxBitsValues[errorCorrectionLevel][version-1]

	b.addPaddingToBytesLen(maxBits / 8)

	blockedData := toBlocks(b.bytes(), errorCorrectionLevel, version)

	correctionBytesLen := correctionBytesValues[errorCorrectionLevel][version-1]

	corr := correctionBytes[correctionBytesLen]

	var corrBlocks [][]byte

	for _, block := range blockedData {
		corrData := correction(block, corr)
		corrBlocks = append(corrBlocks, corrData)
	}

	mergedBlocks := mergeBlocks(blockedData)
	mergedCorr := mergeBlocks(corrBlocks)

	result := append(mergedBlocks, mergedCorr...)

	modulesSize := 21
	if version > 1 {
		alignment := alignmentPositions[version-1]
		modulesSize = alignment[len(alignment)-1] + 7
	}

	bl := newBuilder(modulesSize)

	data := bl.build(result, maskID, errorCorrectionLevel, version)

	img, errRender := renderer.Render(data)
	if errRender != nil {
		return fmt.Errorf("render error: %w", errRender)
	}

	w := bytes.NewBuffer(nil)

	errEncode := png.Encode(w, img)
	if errEncode != nil {
		return fmt.Errorf("encode to png error: %w", errEncode)
	}

	return os.WriteFile(destFile, w.Bytes(), 0644)
}
