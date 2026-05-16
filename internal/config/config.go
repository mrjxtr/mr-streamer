// Package config
package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Platforms struct {
	FB bool
	YT bool
	TW bool
}

type Config struct {
	FBKey string
	YTKey string
	TWKey string

	MRKey string
	PORT  string
}

func LoadConfig(p *Platforms) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn("error loading .env", "error", err)
	}

	cfg := &Config{
		FBKey: os.Getenv("FB_STREAM_KEY"),
		YTKey: os.Getenv("YT_STREAM_KEY"),
		TWKey: os.Getenv("TW_STREAM_KEY"),

		MRKey: os.Getenv("MR_STREAM_KEY"),
		PORT:  os.Getenv("PORT"),
	}
	if err := cfg.validate(p); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate(p *Platforms) error {
	fields := []struct {
		name     string
		val      string
		required bool
	}{
		{"FB_STREAM_KEY", c.FBKey, p.FB},
		{"YT_STREAM_KEY", c.YTKey, p.YT},
		{"TW_STREAM_KEY", c.TWKey, p.TW},
		{"MR_STREAM_KEY", c.MRKey, true},
		{"PORT", c.PORT, true},
	}

	var missing []string
	for _, f := range fields {
		if f.required && f.val == "" {
			missing = append(missing, f.name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing env vars: %s", strings.Join(missing, ", "))
	}

	return nil
}
