package assets

import (
	"bytes"
	"image"
	"image/png"
	"os"
)

const (
	WELLZ_ICON = "./assets/favicon-32x32.png"
)

type Icon struct {
	icon string
}

func New(icon string) *Icon {
	return &Icon{
		icon: icon,
	}
}
func (i *Icon) Get() ([]byte, error) {
	img, err := i.getImageFromFilePath(i.icon)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (i *Icon) getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func t() {

}
