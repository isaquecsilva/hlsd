package m3u8

import (
	"bufio"
	"io"
	"regexp"
	"time"

	"github.com/isaquecsilva/hlsd/config"
	"github.com/isaquecsilva/hlsd/m3u8/uribuilder"
)

var timeRegex = regexp.MustCompile(`\d+\.\d+`)

type HlsParser struct {
	endpoint   string
	r          io.Reader
	uriChan    chan string
	validateTS ParseValidationFunc
	start, end int64
}

func NewHlsParser(
	r io.Reader,
	endpoint string,
	validation ParseValidationFunc,
	start, end int64,
) HlsParser {
	return HlsParser{
		r:          r,
		validateTS: validation,
		endpoint:   endpoint,
		uriChan:    make(chan string),
		start:      start,
		end:        end,
	}
}

func (hp HlsParser) UriChan() chan string {
	return hp.uriChan
}

func (hp HlsParser) Parse() error {
	defer close(hp.uriChan)

	var secondsSum int64

	for line := range hp.lines() {
		// adding the elapsed time to secondsSum
		if timestamp, isTime := hp.isStringTime(line); isTime {
			secondsSum = secondsSum + hp.textToSeconds(timestamp)
		}

		skip, end := hp.timeSkipper(secondsSum)

		if skip || end || !hp.validateTS(line) {
			continue
		}

		uri := uribuilder.Build(hp.endpoint, line)
		if uri == nil {
			continue
		}

		hp.uriChan <- uri.String()
	}

	return nil
}

func (hp HlsParser) lines() chan string {
	linesChan := make(chan string)

	go func() {
		scanner := bufio.NewScanner(hp.r)

		for scanner.Scan() {
			line := scanner.Text()
			linesChan <- line
		}

		close(linesChan)
	}()

	return linesChan
}

func (hp HlsParser) isStringTime(text string) (string, bool) {
	match := timeRegex.FindString(text)
	if match != "" {
		return match, true
	}
	return match, false
}

func (hp HlsParser) textToSeconds(sec string) int64 {
	d, err := time.ParseDuration(sec + "s")

	if err != nil {
		config.Logger().Error(err.Error(), "op", "text-to-seconds")
	}
	return int64(d.Seconds())
}

func (hp HlsParser) timeSkipper(secondsSum int64) (skip, end bool) {
	if secondsSum < hp.start && hp.start > 0 {
		skip = true
	}
	if secondsSum >= hp.end && hp.end > 0 {
		end = true
	}

	return
}
