/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package verify_test

import (
	"github.com/jarednogo/board/backend/verify"
	"testing"
)

func TestHash1(t *testing.T) {
	hash := "$2a$10$O5gxBye06fldX9Al9Upc8.nYkE33KWMTfxYF/sMt5TgEs66Vcg0JK"
	password := "deadbeef"

	if !verify.CorrectPassword(password, hash) {
		t.Errorf("error checking password")
	}
}

func TestVerify1(t *testing.T) {
	msg := []byte("hello world")
	mac := []byte{215, 112, 146, 4, 167, 84, 72, 161, 121, 43, 198, 224, 171, 180, 202, 215, 68, 15, 1, 60, 30, 190, 6, 42, 74, 126, 170, 101, 59, 165, 152, 36}
	key := []byte{0xde, 0xad, 0xbe, 0xef}
	if !verify.VerifyWithKey(msg, mac, key) {
		t.Errorf("verify failed")
	}
}
