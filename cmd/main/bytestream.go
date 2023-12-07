package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jackinthebox52/bytestream/internal/stream"
)

var DEBUG = false

func postBstreams(c *gin.Context) {
	new_stream := stream.ByteStream{}

	if err := c.BindJSON(&new_stream); err != nil {
		c.Status(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	if _, err := stream.CreateStream(new_stream); err != nil {
		c.Status(http.StatusBadRequest)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"uuid":   new_stream.UUID})
}

func getIndex(c *gin.Context) {
	c.HTML(200, "index.tmpl", gin.H{
		"title":   "bytestream - home",
		"streams": stream.STREAMS,
	})
}

func getPlayer(c *gin.Context) {
	if queryParam, ok := c.GetQuery("id"); ok {
		if s, err := stream.GetStreamByUUID(queryParam); err == nil {
			c.HTML(200, "player.tmpl", gin.H{
				"title":     "bsplayer - " + s.StreamName,
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
	s1 := stream.ByteStream{
		StreamURL:      "https://ed-c003.edgking.me/plyvivo/705in0ji3uje8u209ena/chunklist.m3u8",
		StreamName:     "NFL RedZone",
		StreamReferrer: "https://www.niaomea.me/",
		AddedTime:      time.Now(),
		UUID:           stream.GenerateUniqueUUID(),
	}
	s2 := stream.ByteStream{
		StreamURL:      "https://ed-c003.edgking.me/plyvivo/504isoyida6apeloji90/media-uzbcaycgs_5347.ts",
		StreamName:     "ESPN NBA",
		StreamReferrer: "https://www.niaomea.me/",
		AddedTime:      time.Now(),
		UUID:           stream.GenerateUniqueUUID(),
	}
	stream.STREAMS = append(stream.STREAMS, s1, s2)
}

func main() {
	instantiateStreamList() //TODO remove
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "-d" {
			DEBUG = true
		}
	}
	r := gin.Default()

	r.Static("/hls", "./hls")
	r.LoadHTMLGlob("web/templates/*.tmpl")

	r.GET("/", getIndex)
	r.GET("/player", getPlayer)
	r.POST("/bstreams", postBstreams)

	r.Run(":8081") // listen and serve on 0.0.0.0:8081
}
