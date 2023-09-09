package main

import (
	"video-streaming/pkg/streamer"
	"video-streaming/pkg/uploader"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/upload", uploader.Upload)

	r.GET("/play/:video_id/:playlist", streamer.Stream)

	r.Run(":8080")
}
