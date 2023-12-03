package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

type stream struct {
	StreamURL      string    `json:"streamurl"`
	StreamName     string    `json:"streamname"`
	StreamReferrer string    `json:"streamreferrer"`
	AddedTime      time.Time `json:"-"` //JSON ignores this field
}

var SINGLE_STREAM stream

var STREAMS = []stream{}

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

func addStream(c *gin.Context) {
	new_stream := stream{}
	new_stream.AddedTime = time.Now()
	if err := c.BindJSON(&new_stream); err != nil {
		c.Status(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	if err := validateStream(new_stream); err != nil {
		c.Status(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	STREAMS = append(STREAMS, new_stream) //Add the stream to the list
	c.Status(http.StatusCreated)
}

func getIndex(c *gin.Context) {
	c.HTML(200, "index.tmpl", gin.H{
		"title": "bytestream - home",
	})
}

func getPlayer(c *gin.Context) {
	c.HTML(200, "player.tmpl", gin.H{
		"title": "bytestream - player",
	})
}

func instantiateStreamList() {
	s1 := stream{
		StreamURL:      "https://www.youtube.com/watch?v=9E9hHkyZ8Lg",
		StreamName:     "UFC 257",
		StreamReferrer: "https://www.youtube.com/",
		AddedTime:      time.Now(),
	}
	s2 := stream{
		StreamURL:      "https://www.youtube.com/watch?v=9E9hHkyZ8Lg",
		StreamName:     "UFC 257",
		StreamReferrer: "https://www.youtube.com/",
		AddedTime:      time.Now(),
	}
	s3 := stream{
		StreamURL:      "https://www.youtube.com/watch?v=9E9hHkyZ8Lg",
		StreamName:     "UFC 257",
		StreamReferrer: "https://www.youtube.com/",
		AddedTime:      time.Now(),
	}
	STREAMS = append(STREAMS, s1, s2, s3)
}

func main() {
	instantiateStreamList() //TODO remove
	r := gin.Default()

	r.Static("/hls", "./hls")
	r.LoadHTMLGlob("web/*.html")

	r.GET("/", getIndex)
	r.GET("/player", getPlayer)
	r.POST("/addstream", addStream)

	r.Run(":8081") // listen and serve on 0.0.0.0:8081
}
