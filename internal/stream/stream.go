package stream

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"time"
)

var STREAMS = []ByteStream{}

type ByteStream struct {
	StreamURL      string    `json:"streamurl"`
	StreamName     string    `json:"streamname"`
	StreamReferrer string    `json:"streamreferrer"`
	AddedTime      time.Time `json:"-"` //JSON ignores this field
	UUID           string    `json:"-"`
}

// GenerateUUID generates a random 6 character string of alphanumeric characters, and checks if it is unique in the GLOBAL streams list
func GenerateUUID() string {
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
	fmt.Println(s.StreamURL)
	fmt.Println(s.StreamReferrer)
	fmt.Println(s.StreamName)
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

// CreateStream consumes a bytestream object, adds it to the GLOBAL streams lsit, and spawns a goroutine to run ffmpegd
// TODO add channel support for goroutine, add database support
func CreateStream(s ByteStream) (ByteStream, error) {
	if err := validateStream(s); err != nil {
		return ByteStream{}, err
	}
	s.AddedTime = time.Now()
	s.UUID = GenerateUUID()
	if s.StreamName == "" {
		s.StreamName = "Untitled Stream - " + s.UUID
	}
	return s, nil
}
