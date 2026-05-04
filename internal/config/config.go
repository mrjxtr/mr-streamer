package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	FBKey string
	YTKey string
	TWKey string

	MRKey string
	PORT  string
}

func LoadConfig() (*Config, error) {
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
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	fields := []struct{ name, val string }{
		{"FB_STREAM_KEY", c.FBKey},
		{"YT_STREAM_KEY", c.YTKey},
		{"TW_STREAM_KEY", c.TWKey},
		{"MR_STREAM_KEY", c.MRKey},
		{"PORT", c.PORT},
	}
	var missing []string
	for _, f := range fields {
		if f.val == "" {
			missing = append(missing, f.name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing env vars: %s", strings.Join(missing, ", "))
	}

	return nil
}
