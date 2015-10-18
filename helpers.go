/*
Functions in this file should be moved out once a sensible size of
helpers that fit into the same group have accumulated here
*/
package main

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"log"
)

// Marshal is used to easily return JSON
// to the end user without handling errors
// every time in the controller.
func Marshal(i interface{}) []byte {
	response, err := json.Marshal(&i)
	if err != nil {
		log.Fatal(err)
	}

	return response
}

// MakeTarFile take a string and turns it into a tar
// https://golang.org/src/archive/tar/example_test.go
// Most of this is stolen from here. This is specific
// to creating a tar file to be sent to a docker client
func MakeTarFile(content string) (*bytes.Buffer, error) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new tar archive.
	tw := tar.NewWriter(buf)

	// Add some files to the archive.
	var files = []struct {
		Name, Body string
	}{
		{"Dockerfile", content},
	}
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0755,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			return nil, err
		}
	}
	// Make sure to check the error on Close.
	if err := tw.Close(); err != nil {
		return nil, err
	}
	return buf, nil
}
