// License: AGPL-3.0-only
// (c) 2024 Dakota Walsh <kota@nilsu.org>
package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	url := flag.String(
		"url",
		"https://paste.nilsu.org/upload",
		"Upload URL for pastebin server",
	)
	flag.Parse()
	args := flag.Args()

	for _, arg := range args {
		uploadedURL, err := paste(arg, *url)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(uploadedURL)
		}
	}
}

func paste(path, url string) (string, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if stat.IsDir() {
		path, err = zipDir(path)
		if err != nil {
			return "", err
		}
		defer os.Remove(path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	form, err := w.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return "", err
	}

	_, err = form.Write(content)
	if err != nil {
		return "", err
	}

	err = w.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	return res.Request.URL.String(), nil
}

func zipDir(path string) (string, error) {
	zipFile, err := os.CreateTemp("", "*.zip")
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		f, err := zipWriter.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	err = filepath.Walk(path, walker)
	return zipFile.Name(), err
}
