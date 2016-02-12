package goquery

import (
	"compress/gzip"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func NewDocumentFromUrl(cli *http.Client, url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return NewDocumentFromRequest(cli, req)
}

func NewDocumentFromRequest(cli *http.Client, req *http.Request) (*goquery.Document, error) {
	req.Header.Add("Accept-Encoding", "gzip")
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	return goquery.NewDocumentFromReader(reader)
}
