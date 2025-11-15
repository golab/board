/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type mode string

const (
	ModeProd mode = "prod"
	ModeTest mode = "test"
)

type Config struct {
	Mode   mode         `yaml:"mode"`
	Server serverConfig `yaml:"server"`
	Twitch twitchConfig `yaml:"twitch"`
	DB     dbConfig     `yaml:"db"`
}

func (c *Config) Redact() {
	c.Server.redact()
	c.Twitch.redact()
	c.DB.redact()
}

type serverConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	URL  string `yaml:"url"`
}

func (c *serverConfig) redact() {}

type twitchConfig struct {
	ClientID string `yaml:"client_id"`
	Secret   string `yaml:"secret"`
	BotID    string `yaml:"bot_id"`
}

func (c *twitchConfig) redact() {
	if c.ClientID != "" {
		c.ClientID = "***"
	}
	if c.Secret != "" {
		c.Secret = "***"
	}
	if c.BotID != "" {
		c.BotID = "***"
	}
}

type dbConfigType string

const (
	DBConfigTypeSqlite dbConfigType = "sqlite"
	DBConfigTypeMemory dbConfigType = "memory"
)

type dbConfig struct {
	Type dbConfigType `json:"type"`
	Path string       `json:"path"`
}

func (c *dbConfig) redact() {}

func New(fname string) (*Config, error) {
	data, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	// start with base default values
	cfg := Default()

	// the unmarshal step will overwrite values in the default
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if cfg.DB.Type == DBConfigTypeMemory {
		cfg.DB.Path = ""
	}

	return cfg, nil
}

func Default() *Config {
	s := serverConfig{
		Host: "localhost",
		Port: 8080,
		URL:  "http://localhost:8080",
	}
	db := dbConfig{
		Type: DBConfigTypeSqlite,
		Path: defaultSqlitePath(),
	}
	cfg := &Config{
		Mode:   ModeProd,
		Server: s,
		DB:     db,
	}
	return cfg
}

func Test() *Config {
	c := Default()
	c.Mode = ModeTest
	db := dbConfig{
		Type: DBConfigTypeMemory,
	}
	c.DB = db
	return c
}

func defaultSqlitePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	dbPath := filepath.Join(home, ".config", "tripleko", "board.db")
	return dbPath
}
