# mr-streamer 📡

Stream once, broadcast everywhere. A tiny Go supervisor around `ffmpeg` that takes a single RTMP input and fans it out to Facebook, YouTube, and Twitch at the same time.

## Requirements

- Go 1.26+
- `ffmpeg` on your `$PATH` (built with the `tee` muxer, which is the default)

## Setup

1. Copy the example env and fill in your keys:

   ```sh
   cp .example.env .env
   ```

2. Edit `.env`:

   ```env
   FB_STREAM_KEY=your-facebook-key
   YT_STREAM_KEY=your-youtube-key
   TW_STREAM_KEY=your-twitch-key

   MR_STREAM_KEY=any-secret-of-your-choice

   PORT=1935
   ```

   `MR_STREAM_KEY` is any string you make up. It's the local key your streaming software (OBS, etc.) uses to connect to mr-streamer, so treat it like a password. 🔑

## Run

```sh
go run .
```

Or build it:

```sh
go build -o mr-streamer
./mr-streamer
```

## Stream to it

Point OBS (or any RTMP source) at:

```
rtmp://localhost:<PORT>/live/<MR_STREAM_KEY>
```

That single feed gets relayed live to all three platforms. 🎥

## Stop it

`Ctrl+C` once for a graceful shutdown (10s grace before force-kill). Hit it again if you're impatient.

## License

MIT. See [LICENSE](./LICENSE).
