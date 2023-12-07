package ingest

import (
	"fmt"
	"time"

	"github.com/jackinthebox52/bytestream/internal/stream"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FFmpegD struct {
	bs_obj    stream.ByteStream
	proc_id   int
	last_seen time.Time
}

var FFmpegDs = []FFmpegD{}

// IngestHLS ingests a stream using ffmpeg. Spawns a new ffmpeg process and appends it to the global FFmpegDs list. Returns an FFmpegD object
func IngestHLS(bs stream.ByteStream) {
	ref := bs.StreamReferrer
	origin := bs.StreamReferrer[:len(ref)-1]
	err := ffmpeg.Input(bs.StreamURL, ffmpeg.KwArgs{
		"user_agent":          "Mozilla/5.0 (Windows NT 10.0; rv:78.0) Gecko/20100101 Firefox/78.0",
		"headers":             fmt.Sprintf(`Referer: %v\r\nOrigin: %v\r\n`, ref, origin),
		"reconnect":           1,
		"reconnect_at_eof":    1,
		"reconnect_streamed":  1,
		"reconnect_delay_max": 1}).
		Output("./hls/"+bs.UUID+".m3u8", ffmpeg.KwArgs{
			"vcodec":        "copy",
			"hls_time":      10,
			"hls_list_size": 0,
			"start_number":  0,
		}).OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ffmpegd completed for stream " + bs.StreamName + " with UUID " + bs.UUID)

}
