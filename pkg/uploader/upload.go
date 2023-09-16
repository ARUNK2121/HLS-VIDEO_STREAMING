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

	response := make(chan string, 2)
	errChan := make(chan error, 2)

	for _, v := range files {
		wg.Add(1)
		//new go routines
		go ProcessAndUploadFile(&wg, v, response, errChan)
	}

	wg.Wait()
	close(response)

	select {
	case err := <-errChan:
		c.JSON(http.StatusCreated, gin.H{
			"message": err,
		})
		return
	default:
		a := <-response
		b := <-response
		c.JSON(http.StatusCreated, gin.H{
			"message":   "successfully uploaded file to server",
			"video_id":  a,
			"video2_id": b,
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

	fmt.Println("cpas completed and returned nil")
	return nil
}

func ProcessAndUploadFile(wg *sync.WaitGroup, fileheader *multipart.FileHeader, response chan string, errChan chan error) {
	defer wg.Done()
	uuid := uuid.New()
	fileName := uuid.String()
	folderPath := storageLocation + "/" + fileName
	filePath := storageLocation + "/" + fileName + "/" + "video.mp4"
	fmt.Println("1")

	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		errChan <- err
		return
	}
	fmt.Println(2)
	newFile, err := os.Create(filePath)
	if err != nil {
		errChan <- err
		return
	}

	defer newFile.Close()

	src, err := fileheader.Open()
	if err != nil {
		errChan <- err
		return
	}
	fmt.Println(3)
	_, err = io.Copy(newFile, src)
	if err != nil {
		errChan <- err
		return
	}
	fmt.Println(4)
	err = CreatePlaylistAndSegments(filePath, folderPath)
	if err != nil {
		errChan <- err
		return
	}
	fmt.Println("5")
	r := uuid.String()
	fmt.Println("uuid", r)

	response <- r

	fmt.Println("6")

}
