package ingest

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"github.com/jackinthebox52/bytestream/internal"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FFmpegD struct {
	Bstream  ByteStream
	ProcId   *int
	LastSeen time.Time
}

var FFmpegDs = []FFmpegD{}

// Creates a new FFmpegD object and appends it to the global FFmpegDs list
func createFFmd(bs ByteStream) (FFmpegD, error) {
	ffmd := FFmpegD{bs, nil, time.Now()}
	ff, err := IngestHLS_Binary(ffmd)
	if err != nil {
		return FFmpegD{}, err
	}
	FFmpegDs = append(FFmpegDs, ff)
	return ff, nil
}

// Removes an FFmpegD object from the global FFmpegDs list, and kills the process
func removeFFmpegD(uuid string) error { //TODO refactor
	for i, d := range FFmpegDs {
		if d.Bstream.UUID == uuid {
			fmt.Printf("Removing FFmpegD with UUID %v\n", uuid)
			FFmpegDs = append(FFmpegDs[:i], FFmpegDs[i+1:]...)
			if d.ProcId != nil {
				fmt.Printf("Killing process %v\n", *d.ProcId)
				err := exec.Command("kill -9", fmt.Sprintf("%v", *d.ProcId)).Run()
				if err != nil {
					return err
				}
				return nil
			} else {
				fmt.Println("Process ID is nil")
				return nil
			}
		}
	}
	return fmt.Errorf("FFmpegD with UUID %v not found", uuid)
}

func GetFFmpegDByUUID(uuid string) (FFmpegD, error) {
	for _, d := range FFmpegDs {
		if d.Bstream.UUID == uuid {
			return d, nil
		}
	}
	return FFmpegD{}, fmt.Errorf("FFmpegD with UUID %v not found", uuid)
}

func IngestHLS_Binary(ffmd FFmpegD) (FFmpegD, error) {
	bs := ffmd.Bstream
	ref := bs.StreamReferrer
	cmd := exec.Command("script/ffmpegd.sh", bs.StreamURL, bs.UUID, ref)
	var outb bytes.Buffer
	cmd.Stdout = &outb
	//cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
		return FFmpegD{}, err
	}
	in, err := internal.LoadPIDFile(bs.UUID)
	if err != nil {
		fmt.Println(err)
		return FFmpegD{}, err
	}
	ffmd.ProcId = &in
	fmt.Printf("spawned ffmpegd for stream %v with UUID %v, process ID:\n", bs.StreamName, bs.UUID)
	return ffmd, nil
}

/* IngestHLS ingests a stream using ffmpeg. Spawns a new ffmpeg process and appends it to the global FFmpegDs list. Returns an FFmpegD object //NOT IMPLEMENTED ASYNCHRONOUSLY
func IngestHLS_Library(bs ByteStream) {
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
		}).OverWriteOutput().ErrorToStdOut().RunLinux()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ffmpegd completed for stream " + bs.StreamName + " with UUID " + bs.UUID)
}*/

func CleanOldFFmpegDs(hours int, ffmpegds []FFmpegD) []FFmpegD {
	var newFFmpegds []FFmpegD
	for _, d := range ffmpegds {
		if time.Since(d.LastSeen).Hours() < float64(hours) {
			newFFmpegds = append(newFFmpegds, d)
		}
	}
	return newFFmpegds
}

// ArchiveHLS archives a stream using ffmpeg. Takes a ByteStream object as an argument. //ASYNC
func ArchiveHLS(bs ByteStream) {
	go func() {
		err := ffmpeg.Input("./streams/hls/"+bs.UUID+".m3u8", ffmpeg.KwArgs{
			"vcodec":        "copy",
			"hls_time":      10,
			"hls_list_size": 0,
			"start_number":  0,
		}).Output("./hls/"+bs.UUID+".mp4", ffmpeg.KwArgs{
			"vcodec": "copy",
		}).OverWriteOutput().ErrorToStdOut().RunLinux()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Archived stream " + bs.StreamName + " with UUID " + bs.UUID)

	}()
	return
}
