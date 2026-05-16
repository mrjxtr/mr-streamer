package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mrjxtr/mr-streamer/internal/config"
)

const (
	flvOpts = "[f=flv:onfail=ignore]"

	fburl = "rtmps://live-api-s.facebook.com:443/rtmp/"
	yturl = "rtmps://a.rtmps.youtube.com/live2/"
	twurl = "rtmp://live.twitch.tv/app/"

	forceKillDelay = 10 * time.Second
)

func main() {
	fb := flag.Bool("fb", false, "stream to Facebook")
	yt := flag.Bool("yt", false, "stream to YouTube")
	tw := flag.Bool("tw", false, "stream to Twitch")

	flag.Parse()
	if !*fb && !*yt && !*tw {
		*fb, *yt, *tw = true, true, true
	}

	p := &config.Platforms{
		FB: *fb,
		YT: *yt,
		TW: *tw,
	}

	slog.Info("loading config")
	cfg, err := config.LoadConfig(p)
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

	tee := buildTee(cfg, p)

	cmd := exec.CommandContext(
		ctx, "ffmpeg",
		"-hide_banner",
		"-listen", "1",
		"-i", fmt.Sprintf("rtmp://localhost:%s/live/%s", cfg.PORT, cfg.MRKey),
		"-c", "copy",
		"-map", "0",
		"-f", "tee",
		tee,
	)

	cmd.Cancel = func() error { return cmd.Process.Signal(syscall.SIGTERM) }
	cmd.WaitDelay = forceKillDelay
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	slog.Info("starting mr-streamer")
	if err := cmd.Run(); err != nil && !errors.Is(ctx.Err(), context.Canceled) {
		slog.Error("ffmpeg exited", "error", err)
		if code := cmd.ProcessState.ExitCode(); code != -1 {
			os.Exit(code)
		}
		os.Exit(1)
	}
}

func buildTee(cfg *config.Config, p *config.Platforms) string {
	var tee []string

	if p.FB {
		tee = append(tee, flvOpts+fburl+cfg.FBKey)
	}
	if p.YT {
		tee = append(tee, flvOpts+yturl+cfg.YTKey)
	}
	if p.TW {
		tee = append(tee, flvOpts+twurl+cfg.TWKey)
	}

	return strings.Join(tee, "|")
}
