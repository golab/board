/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/pkg/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.Default()

	assert.Equal(t, cfg.Mode, config.ModeProd)
}

func TestTestConfig(t *testing.T) {
	cfg := config.Test()

	assert.Equal(t, cfg.Mode, config.ModeTest)
}

func TestNewConfig(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config.yaml")
	data := `
mode: prod
server:
  host: localhost
  port: 10001
  url: http://localhost:10001
db:
  type: sqlite
  path: /root/config/board.db
twitch:
  client_id: abc
  secret: 123
  bot_id: xyz
`
	err := os.WriteFile(path, []byte(data), 0o644)
	assert.NoError(t, err)

	cfg, err := config.New(path)
	assert.NoError(t, err)
	assert.Equal(t, cfg.Mode, config.ModeProd)
	assert.Equal(t, cfg.Server.Host, "localhost")
	assert.Equal(t, cfg.Server.Port, 10001)
	assert.Equal(t, cfg.Server.URL, "http://localhost:10001")
	assert.Equal(t, cfg.DB.Type, config.DBConfigTypeSqlite)
	assert.Equal(t, cfg.DB.Path, "/root/config/board.db")
	assert.Equal(t, cfg.Twitch.ClientID, "abc")
	assert.Equal(t, cfg.Twitch.Secret, "123")
	assert.Equal(t, cfg.Twitch.BotID, "xyz")

	cfg.Redact()

	assert.Equal(t, cfg.Twitch.ClientID, "***")
	assert.Equal(t, cfg.Twitch.Secret, "***")
	assert.Equal(t, cfg.Twitch.BotID, "***")
}
