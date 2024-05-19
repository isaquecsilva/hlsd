package main

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"
)

type DownloaderInterface interface {
	Enqueue(int, string) error
	QueueFull() bool
	Flush() error
}

type Downloader struct {
	queueSize int
	maxQueue  int
	mu        *sync.Mutex
	streams   map[int][]byte
	w         io.Writer
}

func NewDownloader(queueSize int, w io.Writer) *Downloader {
	downloader := &Downloader{
		maxQueue: queueSize,
		mu:       new(sync.Mutex),
		streams:  make(map[int][]byte),
		w:        w,
	}

	downloader.initDownloaderStreamKeys()
	return downloader

}

func (d *Downloader) initDownloaderStreamKeys() {
	for i := range d.maxQueue {
		d.streams[i] = nil
	}
}

func (d *Downloader) QueueFull() bool {
	return d.queueSize >= d.maxQueue
}

func (d *Downloader) QueueSizeOperation(v int) {
	d.mu.Lock()
	d.queueSize = d.queueSize + v
	d.mu.Unlock()
}

func (d *Downloader) Enqueue(index int, uri string) error {
	response, err := http.Get(uri)

	if err != nil {
		return err
	}

	defer response.Body.Close()
	var buf []byte

	if response.StatusCode != http.StatusOK {
		// second try
		var re = regexp.MustCompile(`\d+`)

		uri = strings.ReplaceAll(path.Dir(uri), ":/", "://") + fmt.Sprintf("/%s-muted.ts", re.FindString(path.Base(uri)))

		response, err = http.Get(uri)

		switch err {
		case nil:
			if response.StatusCode != http.StatusOK {
				return fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}
		default:
			return err
		}
	}

	buf, err = d.download(&uri)

	if err != nil {
		return err
	}

	d.mu.Lock()
	d.streams[index] = buf
	d.mu.Unlock()

	return nil
}

func (d *Downloader) download(uri *string) ([]byte, error) {
	response, err := http.Get(*uri)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	buf, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	} else {
		return buf, nil
	}
}

func (d *Downloader) Flush() error {
	var iterations int = d.queueSize

	for key := range iterations {
		if buf, ok := d.streams[key]; ok {

			if buf == nil {
				// wait until the buffer is ready
				for {
					if buf = d.streams[key]; buf != nil {
						break
					}
					time.Sleep(time.Second)
				}
			}

			_, err := d.w.Write(buf)

			if err != nil {
				return err
			}

			d.streams[key] = nil
			d.QueueSizeOperation(-1)
		} else {
			return fmt.Errorf("downloader.flush: key<%d> not found", key)
		}
	}

	return nil
}
