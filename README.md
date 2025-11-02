# HLS Restreamer

A Docker container that restreams HLS (HTTP Live Streaming) content via a simple API.
Paired with a Firefox extension for quickly adding streams.

## Setup

```bash
git clone [repository-url]
cd hls-restreamer
docker-compose up -d
```

## Firefox Extension

### Installation

1. Open Firefox
2. Navigate to `about:debugging#/runtime/this-firefox`
3. Click "Load Temporary Add-on"
4. Select any file in the `extension/` directory

### Server Configuration

By default, the extension connects to `http://localhost:8080`. To use a remote server:

1. Click the extension icon in the toolbar
2. Click "Settings" (or it will show on first launch)
3. Enter your server URL (e.g., `http://192.168.1.100:8080`)
4. Click "Test Connection" to verify

The server configuration is saved and persists across browser sessions.

## API Usage

### Start a Stream

```bash
curl -X POST http://{SERVER_IP}:8080/api/stream/start \
  -H "Content-Type: application/json" \
  -d '{"name":"NFL", "origin":"https://reliabletv.me", "url":"https://example.com/stream.m3u8"}'
```

Parameters:
- `name`: Stream identifier (alphanumeric, dashes, underscores)
- `origin`: Origin domain for headers (URL format)
- `url`: M3U8 playlist URL to stream

### List Streams

```bash
curl http://{SERVER_IP}:8080/api/streams
```

### Stop a Stream

```bash
curl -X POST http://{SERVER_IP}:8080/api/stream/stop/NFL
```

## Accessing Streams

Access your restreamed content at:
```
http://{SERVER_IP}:8080/stream/STREAM_NAME/index.m3u8
```

## Notes

- Stream segments are stored in `/var/www/streams/` inside the container
- FFmpeg is used with stream copying (no transcoding) to minimize CPU usage
- Streams older than 12 hours are automatically cleaned up