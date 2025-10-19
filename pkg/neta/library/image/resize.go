package image

import (
	"context"
	"fmt"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
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

	// Load image
	img, err := vips.NewImageFromFile(input)
	if err != nil {
		return nil, fmt.Errorf("failed to load image: %w", err)
	}
	defer img.Close()

	// Calculate scale factor
	scale := calculateScale(img, width, height, maintainAspect)

	// Resize
	if err := img.Resize(scale, vips.KernelLanczos3); err != nil {
		return nil, fmt.Errorf("failed to resize image: %w", err)
	}

	// Export
	if err := i.exportImage(img, output, params); err != nil {
		return nil, err
	}

	// Get file size
	fileInfo, _ := os.Stat(output)
	var size int64
	if fileInfo != nil {
		size = fileInfo.Size()
	}

	return map[string]interface{}{
		"path": output,
		"size": size,
		"dimensions": map[string]int{
			"width":  img.Width(),
			"height": img.Height(),
		},
	}, nil
}

// calculateScale determines the scale factor for resizing.
func calculateScale(img *vips.ImageRef, width, height int, maintainAspect bool) float64 {
	if maintainAspect {
		if width > 0 && height == 0 {
			// Scale by width, height will adjust automatically
			return float64(width) / float64(img.Width())
		} else if height > 0 && width == 0 {
			// Scale by height, width will adjust automatically
			return float64(height) / float64(img.Height())
		}
		// Both specified, use smaller scale to fit within bounds
		scaleW := float64(width) / float64(img.Width())
		scaleH := float64(height) / float64(img.Height())
		if scaleW < scaleH {
			return scaleW
		}
		return scaleH
	}

	// No aspect ratio preservation, use width scale
	if width > 0 {
		return float64(width) / float64(img.Width())
	}
	return 1.0 // No scaling
}
