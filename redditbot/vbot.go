package redditbot

import (
	"crypto/rand"
	"fmt"
	"github.com/technicallyty/vidbot/MediaBot"
	"os"
	"strings"
)

// interface guard
var (
	_ MediaBot.Manager = Vbot{}
)

type Vbot struct {
	downloadsDir string
	resourceURL  string
}

//func checkDirExists(path string) bool {
//	dir, err := os.Stat(path)
//	if err != nil {
//		return false
//	}
//	if !dir.IsDir() {
//		return false
//	}
//	return true
//}

func NewVidBot(resourceURL string) Vbot {
	v := Vbot{
		downloadsDir: "",
		resourceURL:  resourceURL,
	}
	v.GetMediaDir()
	return v
}

func (v Vbot) SetResourceURL(url string) bool {
	if isValidRedditLink(url) {
		v.resourceURL = url
		return true
	}
	return false
}

func (v Vbot) Download() (string, string, error) {
	if !v.IsValidURL() {
		return "", "", fmt.Errorf("%v is not a valid reddit URL", v.resourceURL)
	}
	url := removeQueryString(v.resourceURL)
	url = deriveJSONURL(url)

	videoURL, err := fetchMediaURL(url)
	if err != nil {
		return "", "", err
	}
	audioURL, err := deriveAudioURL(videoURL)
	if err != nil {
		return "", "", err
	}

	downloadsDir := v.GetMediaDir()

	uuid := pseudo_uuid()
	resourceName, err := urlToResourceName(v.resourceURL)
	if err != nil {
		return "", "", err
	}

	resourceName += uuid

	//prep the directory for downloads
	if err := os.Mkdir(downloadsDir+"/"+resourceName, os.FileMode(0777)); err != nil {
		return "", "", err
	}

	videoFileName := downloadsDir + "/" + resourceName + "/" + "video.mp4"
	audioFileName := downloadsDir + "/" + resourceName + "/" + "audio.mp4"
	compressedName := downloadsDir + "/" + resourceName + "/" + "compressed.mp4"

	videoChannel := make(chan error)
	audioChannel := make(chan error)

	go func(c chan error) {
		err := MediaBot.DownloadMedia(videoURL, videoFileName)
		c <- err
		close(c)
	}(videoChannel)

	go func(c chan error) {
		err := MediaBot.DownloadMedia(audioURL, audioFileName)
		c <- err
		close(c)
	}(audioChannel)

	select {
	case vidErr := <-videoChannel:
		if vidErr != nil {
			return "", "", vidErr
		}
	case audioErr := <-audioChannel:
		if audioErr != nil {
			path, err := MediaBot.Compress(videoFileName, compressedName)
			if err != nil {
				return "", "", err
			}

			err = os.Remove(videoFileName)
			if err != nil {
				return "", "", err
			}

			return path, compressedName, nil
		}
	}

	combinedName := downloadsDir + "/" + resourceName + "/" + "combined.mp4"
	path, err := MediaBot.Combine(videoFileName, audioFileName, combinedName)
	if err != nil {
		return "", "", err
	}

	path, err = MediaBot.Compress(combinedName, compressedName)
	if err != nil {
		return "", "", err
	}

	return path, resourceName, nil
}

func (v Vbot) ResourceURL() string {
	return v.resourceURL
}

func (v Vbot) IsValidURL() bool {
	return isValidRedditLink(v.resourceURL)
}

func (v Vbot) GetMediaDir() string {
	if v.downloadsDir == "" {
		dir, _ := os.Getwd()
		if strings.Contains(dir, "root") {
			dir := "vidbot/downloads"
			v.downloadsDir = dir
			return dir
		}
		split := strings.Split(dir, "/")
		split[len(split)-1] = "vidbot/downloads"
		dir = strings.Join(split, "/")
		v.downloadsDir = dir
		return dir
	}
	return v.downloadsDir
}

// Note - NOT RFC4122 compliant
func pseudo_uuid() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}
