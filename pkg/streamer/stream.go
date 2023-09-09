package streamer

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"net/http"
)

func Stream(c *gin.Context) {
	// Fetch video id and playlist name from path parameters
	videoID := c.Param("video_id")
	playlist := c.Param("playlist")

	playlistData, err := readPlaylistData(videoID, playlist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to read file from server",
			"error":   err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/vnd.apple.mpegurl")
	c.Header("Content-Disposition", "inline")

	// Write the playlist data to the response body
	c.Writer.Write(playlistData)

}

func readPlaylistData(videoID, playlist string) ([]byte, error) {
	// Construct the playlist file path
	playlistPath := fmt.Sprintf("storage/%s/%s", videoID, playlist)

	// Read the playlist file
	playlistData, err := os.ReadFile(playlistPath)
	if err != nil {
		return nil, err
	}
	return playlistData, nil
}
