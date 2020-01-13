package images

import (
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

var TooSmall = errors.New("data too small")
var TooLarge = errors.New("data too large")

func Test(img image.Image, scope Scope) error {
	r := img.Bounds()
	w := r.Dx()
	h := r.Dy()

	if scope.Min.Width > 0 && w < scope.Min.Width {
		return TooSmall
	}
	if scope.Min.Height > 0 && h < scope.Min.Height {
		return TooSmall
	}
	if scope.Max.Width > 0 && w > scope.Max.Width {
		return TooLarge
	}
	if scope.Max.Height > 0 && h > scope.Max.Height {
		return TooLarge
	}

	return nil
}

func TestFile(path string, scope Scope) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}

	return Test(img, scope)
}
