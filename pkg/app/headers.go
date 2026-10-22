/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package app

import (
	"net/http"
	"strings"
)

// cdn hosts bootstrap and bootstrap-icons (css, js and fonts)
const cdn = "https://cdn.jsdelivr.net"

// hash of board.html's inline <script>0</script>
const boardInlineScript = "'sha256-X+zrZv/IbzjZUnhsbWlsecLbwjndTpG0ZynXOif7V+k='"

// contentSecurityPolicy is the policy for every response. scripts only
// come from this server and the cdn, so injected markup can't run any.
// styles allow 'unsafe-inline' since the pages and scripts set style
// attributes throughout; injected styles can't execute code
// frame-ancestors is left unset so the board can still be embedded
func contentSecurityPolicy(host string) string {
	return strings.Join([]string{
		"default-src 'self'",
		"script-src 'self' " + cdn + " " + boardInlineScript,
		"style-src 'self' 'unsafe-inline' " + cdn,
		"font-src 'self' " + cdn,
		"img-src 'self' data:",
		// 'self' covers same-origin websockets in current browsers; the
		// explicit hosts are for older ones
		"connect-src 'self' ws://" + host + " wss://" + host,
		"object-src 'none'",
		"base-uri 'self'",
		"form-action 'self'",
	}, "; ")
}

// SecurityHeaders sets headers that limit what a browser will do with
// our responses
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Content-Security-Policy", contentSecurityPolicy(r.Host))
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}
