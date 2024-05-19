package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

var (
	endpoint  = flag.String("m3u8", "", "The m3u8 endpoint to use.")
	queueSize = flag.Int("queue", 4, "The size of queue for asynchronous download. Becareful, The greater The size is, more memory it will use.")
	destFile  = flag.String("file", "stream.ts", "The name of The file to store The stream bytes.")
	logs      = flag.Bool("log", false, "Enables logging each resource downloaded.")
	start     = flag.Int("start", 0, "The time in seconds, from where to start downloading The vod.")
	end       = flag.Int("end", 0, "The time in seconds, to stop downloading The vod.")
	version   = flag.Bool("version", false, "Shows current application version.")
)

const currentVersion = "v1.0.0"

func getM3U8Reader(uri string) io.ReadCloser {
	response, err := http.Get(uri)

	if err != nil {
		log.Fatal(err)
	}

	return response.Body
}

func main() {
	flag.Parse()

	if *version {
		fmt.Println(currentVersion)
		return
	}


	if *endpoint == "" {
		log.Fatal("missing m3u8 endpoint.")
	}

	reader := getM3U8Reader(*endpoint)
	defer reader.Close()

	parser := NewHlsParser(reader, path.Dir(*endpoint), DefaultParseValidationFunc, *start, *end)

	file, err := os.Create(*destFile)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	downloader := NewDownloader(*queueSize, file)

	var index int

	go func() {
		if err := parser.Parse(); err != nil {
			log.Fatal(err)
		}
	}()

	for uri := range parser.UriChan() {
		if uri == "" {
			if err := downloader.Flush(); err != nil {
				log.Fatal(err)
			}
			fmt.Println("eof")
			break
		}

		if downloader.QueueFull() {
			if err := downloader.Flush(); err != nil {
				log.Fatal(err)
			}
		}

		go func(i int, u string) {
			if err := downloader.Enqueue(i, u); err != nil {
				log.Fatal(err)
			}
		}(index, uri)

		downloader.QueueSizeOperation(+1)

		index++

		if index >= *queueSize {
			index = 0
		}

		if *logs {
			fmt.Printf("downloading: %s\n", uri)
		}
	}

	time.Sleep(time.Second)
}
