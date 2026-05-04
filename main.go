package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrjxtr/mr-streamer/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("error loading config", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	tee := fmt.Sprintf(
		"[f=flv:onfail=ignore]rtmps://live-api-s.facebook.com:443/rtmp/%s|"+
			"[f=flv:onfail=ignore]rtmps://a.rtmps.youtube.com/live2/%s|"+
			"[f=flv:onfail=ignore]rtmp://live.twitch.tv/app/%s",
		cfg.FBKey, cfg.YTKey, cfg.TWKey,
	)

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-hide_banner",
		"-listen", "1",
		"-i", fmt.Sprintf("rtmp://localhost:%s/live/%s", cfg.PORT, cfg.MRKey),
		"-c", "copy",
		"-map", "0",
		"-f", "tee",
		tee,
	)
	cmd.Cancel = func() error { return cmd.Process.Signal(syscall.SIGTERM) }
	cmd.WaitDelay = 10 * time.Second
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	if err := cmd.Run(); err != nil && !errors.Is(ctx.Err(), context.Canceled) {
		slog.Error("ffmpeg exited", "error", err)
		if code := cmd.ProcessState.ExitCode(); code != -1 {
			os.Exit(code)
		}
		os.Exit(1)
	}
}
