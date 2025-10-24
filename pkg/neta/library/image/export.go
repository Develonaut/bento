package image

import (
	"bytes"
	"fmt"
	stdimage "image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/gen2brain/webp"
)

// exportImage exports image to file in the specified format.
func (i *Image) exportImage(img stdimage.Image, output string, params map[string]interface{}) error {
	quality := getIntParam(params, "quality", 80)
	format := determineFormat(output, params)

	// Check .bentoignore in the target directory
	dir := filepath.Dir(output)
	bentoIgnore, err := loadBentoIgnore(dir)
	if err != nil {
		// If we can't load .bentoignore, log warning but continue
		fmt.Fprintf(os.Stderr, "Warning: failed to load .bentoignore: %v\n", err)
	} else if bentoIgnore.shouldIgnore(output) {
		return fmt.Errorf("file %s is protected by .bentoignore and cannot be overwritten", output)
	}

	data, err := encodeImage(img, format, quality)
	if err != nil {
		return err
	}

	if err := os.WriteFile(output, data, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// determineFormat determines output format from params or file extension.
func determineFormat(output string, params map[string]interface{}) string {
	format := getStringParam(params, "_targetFormat", "")
	if format != "" {
		return format
	}

	ext := strings.ToLower(filepath.Ext(output))
	switch ext {
	case ".webp":
		return "webp"
	case ".jpg", ".jpeg":
		return "jpeg"
	case ".png":
		return "png"
	default:
		return "webp"
	}
}

// encodeImage encodes image in the specified format.
func encodeImage(img stdimage.Image, format string, quality int) ([]byte, error) {
	var buf bytes.Buffer
	var err error

	switch format {
	case "webp":
		err = exportWebp(&buf, img, quality)
	case "jpeg":
		err = exportJpeg(&buf, img, quality)
	case "png":
		err = exportPng(&buf, img)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to export image: %w", err)
	}

	return buf.Bytes(), nil
}

// exportWebp exports image as WebP using gen2brain/webp.
func exportWebp(buf *bytes.Buffer, img stdimage.Image, quality int) error {
	options := webp.Options{
		Lossless: false,
		Quality:  quality,
	}
	return webp.Encode(buf, img, options)
}

// exportJpeg exports image as JPEG using standard library.
func exportJpeg(buf *bytes.Buffer, img stdimage.Image, quality int) error {
	options := &jpeg.Options{
		Quality: quality,
	}
	return jpeg.Encode(buf, img, options)
}

// exportPng exports image as PNG using standard library.
func exportPng(buf *bytes.Buffer, img stdimage.Image) error {
	encoder := &png.Encoder{
		CompressionLevel: png.DefaultCompression,
	}
	return encoder.Encode(buf, img)
}
