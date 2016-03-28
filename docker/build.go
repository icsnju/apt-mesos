package main

import (
	"log"
	"time"
	"bytes"
	"archive/tar"
	"github.com/fsouza/go-dockerclient"
)

func main() {
	client, err := docker.NewClient("http://192.168.99.1:5243")
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	inputbuf, outputbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputbuf)
	tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: 10, ModTime: t, AccessTime: t, ChangeTime: t})
	tr.Write([]byte("FROM base\n"))
	tr.Close()
	opts := docker.BuildImageOptions{
	    Name:         "test",
	    InputStream:  inputbuf,
	    OutputStream: outputbuf,
	}
	if err := client.BuildImage(opts); err != nil {
	    log.Fatal(err)
	}	
}
