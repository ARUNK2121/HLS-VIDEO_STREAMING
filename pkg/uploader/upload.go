package uploader

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	storageLocation = "storage"
)

func Upload(c *gin.Context) {
	file1, err := c.FormFile("video1")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch video file from request",
			"errror":  err.Error(),
		})
		return
	}

	file2, err := c.FormFile("video2")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to fetch video file from request",
			"errror":  err.Error(),
		})
		return
	}

	wg := sync.WaitGroup{}

	files := []*multipart.FileHeader{
		file1, file2,
	}

	for _, v := range files {
		wg.Add(1)
		//new go routines
	}

	wg.Wait()

	uuid := uuid.New()
	fileName := uuid.String()
	folderPath := storageLocation + "/" + fileName
	filePath := storageLocation + "/" + fileName + "/" + "video.mp4"

	err = os.MkdirAll(folderPath, 0755)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to crate directory to store files",
			"error":   err.Error(),
		})
		return
	}

	newFile, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create file to copy video file",
			"error":   err.Error(),
		})
		return
	}

	defer newFile.Close()

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{
			"message": "couldnt open the file",
			"error":   err.Error(),
		})
		return
	}

	written, err := io.Copy(newFile, src)
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{
			"message":       "error while copying file",
			"bytes copyied": written,
			"error":         err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":  "successfully uploaded file to server",
		"video_id": uuid,
	})

	err = CreatePlaylistAndSegments(filePath, folderPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create segments and playlist",
			"error":   err.Error(),
		})
		return
	}

}

func CreatePlaylistAndSegments(filePath string, folderPath string) error {
	segmentDuration := 3
	ffmpegCmd := exec.Command(
		"ffmpeg",
		"-i", filePath,
		"-profile:v", "baseline", // baseline profile is compatible with most devices
		"-level", "3.0",
		"-start_number", "0", // start number segments from 0
		"-hls_time", strconv.Itoa(segmentDuration), //duration of each segment in second
		"-hls_list_size", "0", // keep all segments in the playlist
		"-f", "hls",
		fmt.Sprintf("%s/playlist.m3u8", folderPath),
	)

	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create HLS: %v \nOutput: %s ", err, output)
	}

	return nil
}
