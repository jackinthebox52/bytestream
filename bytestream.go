package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type stream struct {
	StreamURL      string    `json:"streamurl"`
	StreamName     string    `json:"streamname"`
	StreamReferrer string    `json:"streamreferrer"`
	AddedTime      time.Time `json:"-"` //JSON ignores this field
	UUID           string    `json:"-"`
}

var STREAMS = []stream{}

func addStream(c *gin.Context) {
	new_stream := stream{}
	new_stream.AddedTime = time.Now()
	new_stream.UUID = generateUUID()
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
		"title":   "bytestream - home",
		"streams": STREAMS,
	})
}

func getPlayer(c *gin.Context) {
	if queryParam, ok := c.GetQuery("id"); ok {
		if s, err := getStreamByUUID(queryParam); err == nil {
			c.HTML(200, "player.tmpl", gin.H{
				"title":     "bytestream - player",
				"streamurl": s.StreamURL,
			})
		} else {
			c.Status(http.StatusNotFound)
			fmt.Println(err)
		}
	} else {
		c.Status(http.StatusBadRequest)
		fmt.Println("No ID query parameter")
	}
}

func instantiateStreamList() {
	s1 := stream{
		StreamURL:      "https://ed-c002.edgking.me/plyvivo/vesavako80hezi8ofe4i/media-u5e67ox2l_36267.ts",
		StreamName:     "NFL RedZone",
		StreamReferrer: "https://www.niaomea.me/",
		AddedTime:      time.Now(),
		UUID:           generateUUID(),
	}
	s2 := stream{
		StreamURL:      "https://ed-c003.edgking.me/plyvivo/504isoyida6apeloji90/media-uzbcaycgs_5347.ts",
		StreamName:     "ESPN NBA",
		StreamReferrer: "https://www.niaomea.me/",
		AddedTime:      time.Now(),
		UUID:           generateUUID(),
	}
	s3 := stream{
		StreamURL:      "https://www.youtube.com/watch?v=9E9hHkyZ8Lg",
		StreamName:     "UFC 259",
		StreamReferrer: "https://www.niaomea.me/",
		AddedTime:      time.Now(),
		UUID:           generateUUID(),
	}
	STREAMS = append(STREAMS, s1, s2, s3)
}

func main() {
	instantiateStreamList() //TODO remove
	r := gin.Default()

	r.Static("/hls", "./hls")
	r.LoadHTMLGlob("web/templates/*.tmpl")

	r.GET("/", getIndex)
	r.GET("/player", getPlayer)
	r.POST("/addstream", addStream)

	r.Run(":8081") // listen and serve on 0.0.0.0:8081
}
