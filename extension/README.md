# HLS Stream Manager Extension

A Firefox extension that detects M3U8 streams on web pages and sends them to your HLS Restreamer server.

## Features

- Automatically detects HLS (.m3u8) streams during browsing
- Displays a count of detected streams in the toolbar icon
- Manages stream rebroadcasting through a simple interface
- Add, update, and remove streams to/from your restreaming server

## Installation

1. Open Firefox
2. Navigate to `about:debugging#/runtime/this-firefox`
3. Click "Load Temporary Add-on"
4. Select any file in the extension directory

## Usage

1. Browse to a website with HLS video content
2. Click the extension icon in the toolbar
3. View detected streams in the "Detected Streams" tab
4. Click "Add to Server" to start restreaming a detected stream
5. Switch to "Active Streams" to manage your restreams

## Directory Structure

```
extension/
├── manifest.json
├── background.js
├── icons/
│   ├── icon-48.png
│   └── icon-96.png
└── popup/
    ├── popup.html
    ├── popup.css
    └── popup.js
```

## Requirements

- Firefox browser
- Running HLS Restreamer server on http://localhost:8080