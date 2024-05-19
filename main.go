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
	endpoint  = flag.String("m3u8", "", "the m3u8 endpoint to use.")
	queueSize = flag.Int("queue", 4, "the size of queue for asynchronous download. Becareful, the greater the size is, more memory it will use.")
	destFile  = flag.String("file", "stream.ts", "the name of the file to store the stream bytes.")
	logs      = flag.Bool("log", true, "Enables logging each resource downloaded.")
	start     = flag.Int("start", 0, "the time in seconds, from where to start downloading the vod.")
	end       = flag.Int("end", 0, "the time in seconds, to stop downloading the vod.")
)

func getM3U8Reader(uri string) io.ReadCloser {
	response, err := http.Get(uri)

	if err != nil {
		log.Fatal(err)
	}

	return response.Body
}

func main() {
	flag.Parse()

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
