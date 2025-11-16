package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"net/http"
)

//go:embed ui/ui.zip
var staticzipFS []byte

func uiHandler(mux *http.ServeMux) error {
	zipReader, err := zip.NewReader(bytes.NewReader(staticzipFS), int64(len(staticzipFS)))
	if err != nil {
		return err
	}
	mux.Handle("/", http.FileServerFS(zipReader))
	return nil
}
