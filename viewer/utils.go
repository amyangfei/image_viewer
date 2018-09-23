package viewer

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func DetectImageType(data []byte) (string, error) {
	_, fm, err := image.DecodeConfig(bytes.NewReader(data))
	return fm, err
}
