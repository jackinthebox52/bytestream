# HLS Restreamer

A Docker container that restreams HLS (HTTP Live Streaming) content via a simple API.
Paired with a firefox extension for quickly adding streams.

## Setup

```bash
git clone [repository-url]
cd hls-restreamer
docker-compose up -d
```

## API Usage

### Start a Stream

```bash
curl -X POST http://localhost:8080/api/stream/start \
  -H "Content-Type: application/json" \
  -d '{"name":"NFL", "origin":"https://reliabletv.me", "url":"https://example.com/stream.m3u8"}'
```

Parameters:
- `name`: Stream identifier (alphanumeric, dashes, underscores)
- `origin`: Origin domain for headers (URL format)
- `url`: M3U8 playlist URL to stream

### List Streams

```bash
curl http://localhost:8080/api/streams
```

### Stop a Stream

```bash
curl -X POST http://localhost:8080/api/stream/stop/NFL
```

## Accessing Streams

Access your restreamed content at:
```
http://localhost:8080/stream/STREAM_NAME/index.m3u8
```

## Notes

- Stream segments are stored in `/var/www/streams/` inside the container
- FFmpeg is used with stream copying (no transcoding) to minimize CPU usage