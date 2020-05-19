package ucdparser

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FetchURL fetches a file from a remote location.
func FetchURL(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET: %v", err)
	}
	if resp.StatusCode != 200 {
		err := fmt.Errorf("bad GET status for '%s': '%s'", url, resp.Status)
		resp.Body.Close()
		return nil, err
	}
	return resp.Body, nil
}

// FetchFile opens local file
func FetchFile(filename string) (io.ReadCloser, error) {
	f, err := os.Open(filename)
	return f, err
}

// Fetch opens url or local file, depending on location prefix
func Fetch(location string) (io.ReadCloser, error) {
	if strings.HasPrefix(location, "http:") || strings.HasPrefix(location, "https:") {
		return FetchURL(location)
	}
	return FetchFile(location)
}

// FetchCached obtains content from local file (cache)
// or, if file is unavailable, fetches it from the
// remote url. Use forceRemote flag to force remote
// (re)fetch regardless of the existance of the
// local file.
func FetchCached(url, file string, forceRemote bool) (io.ReadCloser, error) {
	if !forceRemote {
		// see if cache is available
		if f, err := os.Open(file); err == nil {
			return f, nil
		}
	}

	// fetch from remote
	r, err := Fetch(url)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// store fetched in cache
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("download failure: %v", err)
	}
	os.MkdirAll(filepath.Dir(file), 0755)
	if err := ioutil.WriteFile(file, b, 0755); err != nil {
		log.Fatalf("could not create file: %v", err)
	}

	return ioutil.NopCloser(bytes.NewReader(b)), nil
}
