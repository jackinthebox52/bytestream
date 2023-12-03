package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type stream struct {
	StreamURL      string `json:"streamURL"`
	StreamName     string `json:"streamname"`
	StreamReferrer string `json:"streamreferrer"`
}

var SINGLE_STREAM stream

func addStream(c *gin.Context) {

	if err := c.BindJSON(&s); err != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getIndex(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "bytestream - home",
	})
}

func getUFC(c *gin.Context) {
	c.HTML(200, "hls.html", gin.H{
		"title": "bytestream - ufc",
	})
}

func main() {
	r := gin.Default()

	r.Static("/hls", "./hls")
	r.Static("/public", "./public")
	r.LoadHTMLGlob("web/*.html")

	r.GET("/", getIndex)
	r.GET("/stream", getUFC)

	r.Run(":8081") // listen and serve on 0.0.0.0:8081
}
