package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"
)

func generateUUID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for {
		b := make([]byte, 6)
		for i := range b {
			n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			b[i] = charset[n.Int64()]
		}
		uuid := string(b)
		if _, err := getStreamByUUID(uuid); err != nil {
			return uuid //UUID is unique if getStreamByUUID returns an error
		}
	}
}

func getStreamByUUID(uuid string) (stream, error) {
	for _, s := range STREAMS {
		if s.UUID == uuid {
			return s, nil
		}
	}
	return stream{}, errors.New("Stream not found")
}

func validateStream(s stream) error {
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
