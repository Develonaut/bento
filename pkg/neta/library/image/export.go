package image

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/davidbyttow/govips/v2/vips"
)

// exportImage exports image to file in the specified format.
func (i *Image) exportImage(img *vips.ImageRef, output string, params map[string]interface{}) error {
	quality := getIntParam(params, "quality", 80)
	format := determineFormat(output, params)

	data, err := encodeImage(img, format, quality)
	if err != nil {
		return err
	}

	// Write to file
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
		return "webp" // Default to webp
	}
}

// encodeImage encodes image in the specified format.
func encodeImage(img *vips.ImageRef, format string, quality int) ([]byte, error) {
	var data []byte
	var err error

	switch format {
	case "webp":
		data, err = exportWebp(img, quality)
	case "jpeg":
		data, err = exportJpeg(img, quality)
	case "png":
		data, err = exportPng(img)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to export image: %w", err)
	}

	return data, nil
}

// exportWebp exports image as WebP.
func exportWebp(img *vips.ImageRef, quality int) ([]byte, error) {
	exportParams := vips.NewWebpExportParams()
	exportParams.Quality = quality
	exportParams.StripMetadata = true
	data, _, err := img.ExportWebp(exportParams)
	return data, err
}

// exportJpeg exports image as JPEG.
func exportJpeg(img *vips.ImageRef, quality int) ([]byte, error) {
	exportParams := vips.NewJpegExportParams()
	exportParams.Quality = quality
	exportParams.StripMetadata = true
	data, _, err := img.ExportJpeg(exportParams)
	return data, err
}

// exportPng exports image as PNG.
func exportPng(img *vips.ImageRef) ([]byte, error) {
	exportParams := vips.NewPngExportParams()
	exportParams.StripMetadata = true
	data, _, err := img.ExportPng(exportParams)
	return data, err
}
