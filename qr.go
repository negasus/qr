package qr

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
)

func MakeImage(src []byte) (image.Image, error) {
	d, e := EncodeBytes(ErrorCorrectionLevelM, 0, src)
	if e != nil {
		return nil, e
	}
	return renderImage(d)
}

func SaveImage(src []byte, filename string) error {
	img, e := MakeImage(src)
	if e != nil {
		return e
	}

	w := bytes.NewBuffer(nil)

	errEncode := png.Encode(w, img)
	if errEncode != nil {
		return errEncode
	}

	return os.WriteFile(filename, w.Bytes(), 0644)
}

func EncodeBytes(errorCorrectionLevel ErrorCorrectionLevel, maskID int, value []byte) ([]byte, error) {
	if errorCorrectionLevel < 0 || errorCorrectionLevel > 3 {
		return nil, fmt.Errorf("errorCorrectionLevel must be between 0 and 3")
	}

	if maskID < 0 || maskID > 7 {
		return nil, fmt.Errorf("maskID must be between 0 and 7")
	}

	b := &bitwise{}
	b.addBytes(value)

	var version, ls int

	version = calcVersion(errorCorrectionLevel, len(value)*8)
	if version < 0 {
		return nil, fmt.Errorf("data is too big")
	}

	for {
		ls = getLengthFieldSize(EncodeTypeByte, version)

		if len(value)*8+ls+4 <= maxBitsValues[errorCorrectionLevel][version-1] {
			break
		}

		version++
		if version > 40 {
			return nil, fmt.Errorf("data is too big")
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

	return data, nil
}
