/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package zip_test

import (
	"testing"

	"github.com/golab/board/internal/zip"
)

func TestZip(t *testing.T) {
	data := []byte{80, 75, 3, 4, 10, 0, 0, 0, 0, 0, 131, 158, 9, 91, 32, 48, 58, 54, 6, 0, 0, 0, 6, 0, 0, 0, 5, 0, 28, 0, 102, 105, 108, 101, 49, 85, 84, 9, 0, 3, 69, 251, 151, 104, 69, 251, 151, 104, 117, 120, 11, 0, 1, 4, 232, 3, 0, 0, 4, 232, 3, 0, 0, 104, 101, 108, 108, 111, 10, 80, 75, 3, 4, 10, 0, 0, 0, 0, 0, 132, 158, 9, 91, 168, 97, 56, 221, 6, 0, 0, 0, 6, 0, 0, 0, 5, 0, 28, 0, 102, 105, 108, 101, 50, 85, 84, 9, 0, 3, 72, 251, 151, 104, 72, 251, 151, 104, 117, 120, 11, 0, 1, 4, 232, 3, 0, 0, 4, 232, 3, 0, 0, 119, 111, 114, 108, 100, 10, 80, 75, 1, 2, 30, 3, 10, 0, 0, 0, 0, 0, 131, 158, 9, 91, 32, 48, 58, 54, 6, 0, 0, 0, 6, 0, 0, 0, 5, 0, 24, 0, 0, 0, 0, 0, 1, 0, 0, 0, 164, 129, 0, 0, 0, 0, 102, 105, 108, 101, 49, 85, 84, 5, 0, 3, 69, 251, 151, 104, 117, 120, 11, 0, 1, 4, 232, 3, 0, 0, 4, 232, 3, 0, 0, 80, 75, 1, 2, 30, 3, 10, 0, 0, 0, 0, 0, 132, 158, 9, 91, 168, 97, 56, 221, 6, 0, 0, 0, 6, 0, 0, 0, 5, 0, 24, 0, 0, 0, 0, 0, 1, 0, 0, 0, 164, 129, 69, 0, 0, 0, 102, 105, 108, 101, 50, 85, 84, 5, 0, 3, 72, 251, 151, 104, 117, 120, 11, 0, 1, 4, 232, 3, 0, 0, 4, 232, 3, 0, 0, 80, 75, 5, 6, 0, 0, 0, 0, 2, 0, 2, 0, 150, 0, 0, 0, 138, 0, 0, 0, 0, 0}
	if !zip.IsZipFile(data) {
		t.Errorf("error detecting zip file")
	}

	files, err := zip.Decompress(data)
	if err != nil {
		t.Error(err)
	}

	if len(files) != 2 {
		t.Errorf("error in number of decompressed files")
	}

	if string(files[0]) != "hello\n" {
		t.Errorf("error decompressing, expected %s, got: %s", "hello", string(files[0]))
	}

	if string(files[1]) != "world\n" {
		t.Errorf("error decompressing, expected %s, got: %s", "world", string(files[1]))
	}

}
