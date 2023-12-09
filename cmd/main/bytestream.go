package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jackinthebox52/bytestream/internal/ingest"
	"github.com/jackinthebox52/bytestream/internal/paths"
)

var DEBUG = false

func postBstreams(c *gin.Context) {
	new_stream := ingest.ByteStream{}

	if err := c.BindJSON(&new_stream); err != nil {
		c.Status(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	new_stream, err := ingest.CreateStream(new_stream)
	if err != nil {
		c.Status(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"uuid":   new_stream.UUID})
}

func getIndex(c *gin.Context) {
	c.HTML(200, "index.tmpl", gin.H{
		"title":   "bytestream - home",
		"streams": ingest.STREAMS,
	})
}

func getPlayer(c *gin.Context) {
	if queryParam, ok := c.GetQuery("id"); ok {
		if s, err := ingest.GetStreamByUUID(queryParam); err == nil {
			c.HTML(200, "player.tmpl", gin.H{
				"title": "bsplayer - " + s.StreamName,
				"UUID":  template.JS(s.UUID),
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

type DeleteRequest struct {
	UUID string `json:"uuid"`
}

func deleteBstreams(c *gin.Context) {
	var request DeleteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	if err := ingest.DeleteStream(request.UUID); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"uuid":   request.UUID,
		})
	} else {
		c.Status(http.StatusNotFound)
		fmt.Println(err)
	}
}

func instantiateStreamList() {
	s1 := ingest.ByteStream{
		StreamURL:      "https://ed-c003.edgking.me/plyvivo/705in0ji3uje8u209ena/chunklist.m3u8",
		StreamName:     "NFL RedZone",
		StreamReferrer: "https://www.niaomea.me/",
		AddedTime:      time.Now(),
		UUID:           ingest.GenerateUniqueUUID(),
	}
	s2 := ingest.ByteStream{
		StreamURL:      "https://ed-c003.edgking.me/plyvivo/504isoyida6apeloji90/media-uzbcaycgs_5347.ts",
		StreamName:     "ESPN NBA",
		StreamReferrer: "https://www.niaomea.me/",
		AddedTime:      time.Now(),
		UUID:           ingest.GenerateUniqueUUID(),
	}
	ingest.STREAMS = append(ingest.STREAMS, s1, s2)
}

func main() {
	//instantiateStreamList() //TODO remove
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "-d" {
			DEBUG = true
		}
	}
	r := gin.Default()

	if hlsDir, err := paths.CompileHlsBase(); err != nil {
		panic(err)
	} else {
		r.Static("/hls", hlsDir)
	}
	r.LoadHTMLGlob("web/templates/*.tmpl")

	r.GET("/", getIndex)
	r.GET("/player", getPlayer)
	r.POST("/bstreams", postBstreams)
	r.POST("/rmstream", deleteBstreams)

	r.Run(":8081") // listen and serve on 0.0.0.0:8081
}
