package MediaBot

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

const (
	app = "ffmpeg"
)

// Combine executes an ffmpeg command on the terminal, combining an audio file and video file
func Combine(file1, file2, resultName string) (string, error) {
	args := []string{"-i", file1, "-i", file2, "-c:v", "copy", "-c:a", "aac", resultName}
	cmd := exec.Command(app, args...)
	err := cmd.Run()
	if err != nil {
		return resultName, fmt.Errorf("combine error: %v\n", err)
	}
	return resultName, nil
}

// Compress a file using libx264 compression algorithm
func Compress(file, out string) (string, error) {
	args := []string{"-i", file, "-vcodec", "libx264", "-crf", "28", out}
	cmd := exec.Command(app, args...)
	err := cmd.Run()
	if err != nil {
		return out, fmt.Errorf("compress error: %v\n", err)
	}
	return out, nil
}

// DownloadMedia downloads the media at the specified URL to the specified path
func DownloadMedia(url, savePath string) error {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		return errors.New("error: could not download audio")
	}

	out, err := os.Create(savePath)
	defer out.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
