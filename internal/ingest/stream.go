package ingest

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"time"
)

var STREAMS = []ByteStream{}

type ByteStream struct {
	StreamURL      string    `json:"streamurl"`      //The URL of the external source stream.
	StreamName     string    `json:"streamname"`     //The name of the stream (arbitrary).
	StreamReferrer string    `json:"streamreferrer"` //The URL of the page that the stream is embedded on, or the proper referrer for the stream.
	AddedTime      time.Time `json:"-"`              //JSON ignores this field
	UUID           string    `json:"-"`
}

// GenerateUniqueUUID generates a random 6 character string of alphanumeric characters, and checks if it is unique in the GLOBAL streams list
func GenerateUniqueUUID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for {
		b := make([]byte, 6)
		for i := range b {
			n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			b[i] = charset[n.Int64()]
		}
		uuid := string(b)
		if _, err := GetStreamByUUID(uuid); err != nil {
			return uuid //UUID is unique if getStreamByUUID returns an error
		}
	}
}

func GetStreamByUUID(uuid string) (ByteStream, error) {
	for _, s := range STREAMS {
		if s.UUID == uuid {
			return s, nil
		}
	}
	return ByteStream{}, errors.New("Stream not found")
}

func validateStream(s ByteStream) error {
	urlRegex := `^(http(s):\/\/.)[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`
	if s.StreamURL == "" || s.StreamReferrer == "" {
		return errors.New("Stream URL or Referrer is empty")
	}

	match, err := regexp.MatchString(urlRegex, s.StreamURL)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("Invalid Stream URL, must be a valid URL starting with http(s)://")
	}

	return nil
}

// CreateStream consumes a bytestream object, adds it to the GLOBAL streams lsit, qeueing it for ingestion
// TODO add channel support for goroutine, add database support
func CreateStream(s ByteStream) (ByteStream, error) {
	if err := validateStream(s); err != nil {
		return ByteStream{}, err
	}
	s.AddedTime = time.Now()
	s.UUID = GenerateUniqueUUID()
	if s.StreamName == "" {
		s.StreamName = "Untitled Stream - " + s.UUID
	}
	initalizeDirectory(s.UUID)
	STREAMS = append(STREAMS, s)
	SpawnDeleteStream()
	return s, nil
}

func initalizeDirectory(uuid string) error {
	dirs := []string{fmt.Sprintf("streams/hls/%s", uuid)}
	for _, d := range dirs {
		if err := os.MkdirAll(d, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func DeleteStream(uuid string) error {
	for i, s := range STREAMS {
		if s.UUID == uuid {
			STREAMS = append(STREAMS[:i], STREAMS[i+1:]...)
			SpawnDeleteStream() //Must be called after adding/deleting stream to update the ffmpegd list
			return nil
		}
	}
	return errors.New("Stream not found")
}

// SpawnDeleteStream iterated through the GLOBAL streams list and finds any streams that need created or deleted. It does this by checking the ByteStreams list against the FFmpegDs list
func SpawnDeleteStream() error {
	for _, s := range FFmpegDs {
		fmt.Printf("Checking stream %v\n", s.Bstream.StreamName)
		bs := s.Bstream
		if _, err := GetStreamByUUID(bs.UUID); err != nil { //If stream is not found in database, but is found in FFmpegDs list, delete it
			fmt.Printf("Stream %v not found in database, deleting\n", bs.StreamName)
			//FFmpegDs = append(FFmpegDs[:i], FFmpegDs[i+1:]...)
			removeFFmpegD(bs.UUID)
		}
	}

	for _, s := range STREAMS {
		if _, err := GetFFmpegDByUUID(s.UUID); err != nil { //If stream is not found in FFmpegDs list, but is found in database, create it
			fmt.Printf("FFmpeg %v not found in FFmpegDs list, creating\n", s.StreamName)
			_, err = createFFmd(s)
			if err != nil {
				return err
			}
			fmt.Printf("Created stream %v\n", s.StreamName)
		}
	}
	return nil
}
