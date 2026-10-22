/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// zip provides simple functions for detecting and decompressing zip files
package zip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
)

// limits on what an uploaded archive can expand to
// (a 1 MB zip can otherwise hold 1 GB of zeros)
const (
	MaxFiles      = 256
	MaxTotalBytes = 8 << 20
)

// IsZipFile checks the first two bytes to see if they are 0x504b
func IsZipFile(data []byte) bool {
	return len(data) > 2 && data[0] == 0x50 && data[1] == 0x4b
}

// Decompress will return an array of files (each file is a []byte)
func Decompress(data []byte) ([][]byte, error) {
	// create a reader from the byte slice
	r := bytes.NewReader(data)

	// open the archive with a zip reader
	zipReader, err := zip.NewReader(r, int64(len(data)))
	if err != nil {
		return nil, err
	}

	if len(zipReader.File) > MaxFiles {
		return nil, fmt.Errorf("too many files in archive (max %d)", MaxFiles)
	}

	files := [][]byte{}
	remaining := int64(MaxTotalBytes)
	for _, file := range zipReader.File {
		rc, err := file.Open()
		if err != nil {
			continue
		}
		// read one byte past the budget to detect going over it
		fData, err := io.ReadAll(io.LimitReader(rc, remaining+1))
		rc.Close() //nolint: errcheck
		if err != nil {
			continue
		}
		remaining -= int64(len(fData))
		if remaining < 0 {
			return nil, fmt.Errorf("archive too large when decompressed (max %d MB)", MaxTotalBytes>>20)
		}
		files = append(files, fData)
	}
	return files, nil
}
