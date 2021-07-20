package redditbot

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var redditRegExp = regexp.MustCompile("https:\\/\\/(www.)?reddit.com\\/r\\/\\w+\\/\\w+\\/\\w+\\/\\w+")

func isValidRedditLink(url string) bool {
	return redditRegExp.MatchString(url)
}

// fetches the json information from the given URL. must conform to Response.go
func fetchMediaURL(url string) (string, error) {
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 3,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "vreddit-getter")
	res, getErr := client.Do(req)
	if getErr != nil {
		return "", getErr
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return "", readErr
	}

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	if len(data) < 1 {
		return "", errors.New("bad URL: no data in response")
	}
	videoURL := data[0].Data.Children[0].Data.SecureMedia.RedditVideo.FallbackURL
	if len(videoURL) < 1 {
		return "", errors.New("no MediaBot found")
	}
	if data[0].Data.Children[0].Data.SecureMedia.RedditVideo.Duration > 90 {
		return "", errors.New("error: ribbot alpha can currently only handle video lengths of 1 min 30s or less")
	}
	return videoURL, nil
}

// takes a `v.reddit` url and transforms it to a URL to use for audio MediaBot requests
func deriveAudioURL(url string) (string, error) {
	simplifiedURL := strings.Split(url, "?")
	baseURL := simplifiedURL[0]
	noDash := strings.Split(baseURL, "DASH_")
	if len(noDash) != 2 {
		return "", errors.New("invalid URL: missing DASH identifier in derived URL")
	}
	audioURL := noDash[0] + "DASH_audio.mp4"
	return audioURL, nil
}

// transforms a URL into a JSON url
func deriveJSONURL(url string) string {
	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
		return url + ".json"
	}
	return url + ".json"
}

// removeQueryString removes the query string from a URL
func removeQueryString(url string) string {
	return strings.Split(url, "?")[0]
}

func urlToResourceName(url string) (string, error) {
	if !isValidRedditLink(url) {
		return "", errors.New("not a valid reddit link")
	}
	url = removeQueryString(url)

	// remove trailing slash
	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}

	split := strings.Split(url, "/")
	if len(split) != 8 {
		return "", errors.New("not a valid reddit link: expected https://reddit.com/r/<sub>/comments/<cid>/<post-title>")
	}

	sub := split[4]
	id := split[6]
	title := split[7]

	return sub + "_" + id + "_" + title, nil
}
