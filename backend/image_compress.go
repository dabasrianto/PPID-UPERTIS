package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

const (
	maxImageWidth   = 1920
	jpegQuality     = 75
	compressMinSize = 200 * 1024 // Only compress files > 200KB
)

// compressImage reads an image file, resizes if > maxWidth, and re-encodes as JPEG quality 75%.
// Replaces the original file. Skips GIFs, PDFs, ICOs, and small files.
func compressImage(filePath string) {
	ext := strings.ToLower(filepath.Ext(filePath))

	// Skip non-compressible formats
	if ext == ".gif" || ext == ".pdf" || ext == ".ico" || ext == ".webp" {
		return
	}

	// Check file size — skip if already small
	info, err := os.Stat(filePath)
	if err != nil || info.Size() < compressMinSize {
		return
	}

	// Open file
	f, err := os.Open(filePath)
	if err != nil {
		return
	}

	// Decode image
	var img image.Image
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(f)
	case ".png":
		img, err = png.Decode(f)
	default:
		f.Close()
		return
	}
	f.Close()
	if err != nil {
		return
	}

	// Resize if too wide
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width > maxImageWidth {
		newWidth := maxImageWidth
		newHeight := int(float64(height) * float64(newWidth) / float64(width))

		resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		draw.BiLinear.Scale(resized, resized.Bounds(), img, bounds, draw.Over, nil)
		img = resized
	}

	// Re-encode as JPEG (always output .jpg for compression)
	// Write to same path but with .jpg extension
	outputPath := filePath
	if ext == ".png" {
		// For PNG: save compressed version, keep same filename but as JPEG
		outputPath = strings.TrimSuffix(filePath, ext) + ".jpg"
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return
	}
	defer out.Close()

	jpeg.Encode(out, img, &jpeg.Options{Quality: jpegQuality})

	// If we converted PNG to JPG, remove original PNG and update path
	if ext == ".png" && outputPath != filePath {
		os.Remove(filePath)
	}

	// Log compression
	newInfo, _ := os.Stat(outputPath)
	if newInfo != nil && info.Size() > 0 {
		saved := info.Size() - newInfo.Size()
		if saved > 0 {
			zlog.Info().
				Str("file", filepath.Base(outputPath)).
				Int64("original", info.Size()).
				Int64("compressed", newInfo.Size()).
				Int64("saved_bytes", saved).
				Msg("Image compressed")
		}
	}
}

// getCompressedFilename returns the final filename after compression
// (PNG becomes JPG, others stay same)
func getCompressedExt(originalExt string) string {
	if originalExt == ".png" {
		return ".jpg"
	}
	return originalExt
}
