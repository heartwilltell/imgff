package imgff

import (
	"errors"
	"fmt"
	"os"
)

// Enumeration of supported image formats.
const (
	JPEG FileFormat = "jpg"
	// PNG represents png image format.
	PNG FileFormat = "png"
	// BMP represents bmp image format.
	BMP FileFormat = "bmp"
	// GIF represents gif image format.
	GIF FileFormat = "gif"
	// AVIF represents avif image format.
	AVIF FileFormat = "avif"
	// WEBP represents webp image format.
	WEBP FileFormat = "webp"
)

// FileFormat represents file format.
type FileFormat string

func (f FileFormat) String() string { return string(f) }

const (
	// ErrUnknownFormat means that format can't be determined.
	ErrUnknownFormat Error = "unknown format"
)

// Error represents package level errors.
// Implements builtin.error interface.
type Error string

func (e Error) Error() string { return string(e) }

// Enumeration of special markers (bytes) in the beginning of
// the file (header) by which we can determine the file format.
const (
	// jpeg header bytes
	jpeg0 = 0xFF
	jpeg1 = 0xD8

	// bmp header bytes
	bmp0 = 0x42
	bmp1 = 0x4D

	// png header bytes
	png0 = 0x89
	png1 = 0x50
	png2 = 0x4E
	png3 = 0x47

	// gif header bytes
	gif0 = 0x47
	gif1 = 0x49
	gif2 = 0x46
	gif3 = 0x38

	// TODO: avif header bytes

	// TODO: webp header bytes
)

// Format tries to determine image file format.
func Format(file *os.File) (FileFormat, error) {
	return format(file)
}

// FormatMust tries to determine image file format.
// In rare cases can panic if error occurs.
func FormatMust(file *os.File) FileFormat {
	f, err := format(file)
	if err != nil && !errors.Is(err, ErrUnknownFormat) {
		panic(err)
	}

	return f
}

// format checks and returns the file format.
func format(file *os.File) (FileFormat, error) {
	b := make([]byte, 4)

	n, readErr := file.ReadAt(b, 0)
	if readErr != nil {
		return "", fmt.Errorf("failed to read bytes from file: %w", readErr)
	}

	if n < 4 {
		return "", fmt.Errorf("failed to determine file format: malformed file")
	}

	if b[0] == jpeg0 && b[1] == jpeg1 {
		return JPEG, nil
	}

	if b[0] == bmp0 && b[1] == bmp1 {
		return BMP, nil
	}

	if b[0] == png0 && b[1] == png1 && b[2] == png2 && b[3] == png3 {
		return PNG, nil
	}

	if b[0] == gif0 && b[1] == gif1 && b[2] == gif2 && b[3] == gif3 {
		return GIF, nil
	}

	return "", ErrUnknownFormat
}
