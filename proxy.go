package go_util

import (
	"io"
	"log"
	"net/http"
)

func ProxyFile(cli *http.Client, w http.ResponseWriter, r *http.Request, dst string, refer string) error {
	{ // get cookie
		resp, err := cli.Get(refer)
		if err != nil {
			return err
		}
		resp.Body.Close()
	}

	req, err := http.NewRequest("GET", dst, nil)
	if len(refer) > 0 {
		req.Header.Add("Referer", refer)
	}
	resp, err := cli.Do(req)
	if err != nil {
		log.Printf("Failed to proxy %v: %v", dst, err)
		return err
	}
	if resp.StatusCode >= 400 {
		log.Printf("Failed to proxy %v: %v", dst, resp.Status)
	}
	return proxyResponse(w, r, resp)
}

func proxyResponse(w http.ResponseWriter, r *http.Request, resp *http.Response) error {
	defer resp.Body.Close()

	header := w.Header()
	for k, v := range resp.Header {
		header[k] = v
	}
	w.WriteHeader(resp.StatusCode)

	_, err := io.Copy(w, resp.Body)
	return err
}
