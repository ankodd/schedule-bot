package downloader

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type FileDownloader struct {
	URL      string
	FileName string
}

func DownloadFile(rawUrl string) (string, error) {
	// Parse URL
	URL, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	// Get file name
	URLSegments := strings.Split(URL.Path, "/")
	filename := URLSegments[len(URLSegments)-1]

	// Create file
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Download file
	resp, err := http.DefaultClient.Do(
		&http.Request{
			URL:    URL,
			Method: http.MethodGet,
		},
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Copying the contents of the file
	size, err := io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	log.Printf("Downloaded a file %s with size %d", filename, size)

	return filename, nil
}
