package viewer

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/satori/go.uuid"
	"github.com/teris-io/shortid"
)

func DetectImageType(data []byte) (string, error) {
	_, fm, err := image.DecodeConfig(bytes.NewReader(data))
	return fm, err
}

func RandomId(sid *shortid.Shortid) string {
	if sid != nil {
		if id, err := sid.Generate(); err == nil {
			return id
		}
	}
	return uuid.NewV4().String()
}
