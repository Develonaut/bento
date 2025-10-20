package image

import (
	"context"
	"fmt"
	stdimage "image"
	"os"

	"github.com/disintegration/imaging"
)

// resize resizes an image.
func (i *Image) resize(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	input, ok := params["input"].(string)
	if !ok {
		return nil, fmt.Errorf("input parameter is required")
	}

	output, ok := params["output"].(string)
	if !ok {
		return nil, fmt.Errorf("output parameter is required")
	}

	width := getIntParam(params, "width", 0)
	height := getIntParam(params, "height", 0)
	maintainAspect := getBoolParam(params, "maintainAspect", true)

	if width == 0 && height == 0 {
		return nil, fmt.Errorf("either width or height must be specified")
	}

	img, err := imaging.Open(input)
	if err != nil {
		return nil, fmt.Errorf("failed to load image: %w", err)
	}

	var resized stdimage.Image
	if maintainAspect {
		resized = imaging.Resize(img, width, height, imaging.Lanczos)
	} else {
		if width == 0 {
			width = img.Bounds().Dx()
		}
		if height == 0 {
			height = img.Bounds().Dy()
		}
		resized = imaging.Resize(img, width, height, imaging.Lanczos)
	}

	if err := i.exportImage(resized, output, params); err != nil {
		return nil, err
	}

	fileInfo, _ := os.Stat(output)
	var size int64
	if fileInfo != nil {
		size = fileInfo.Size()
	}

	return map[string]interface{}{
		"path": output,
		"size": size,
		"dimensions": map[string]int{
			"width":  resized.Bounds().Dx(),
			"height": resized.Bounds().Dy(),
		},
	}, nil
}
