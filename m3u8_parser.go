package main

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type ParseValidationFunc = func(string) bool

// Regexes
var (
	tsRegex   = regexp.MustCompile(`\.ts`)
	timeRegex = regexp.MustCompile(`\d+\.\d+,`)
)

var DefaultParseValidationFunc = func(text string) bool {
	return tsRegex.MatchString(text)
}

type HlsParserInterface interface {
	Parse() error
	UriChan() chan string
}

type HlsParser struct {
	rootEndpoint string
	r            io.Reader
	uriChan      chan string
	validateTS   ParseValidationFunc
	start, end   int
}

func NewHlsParser(r io.Reader, endpoint string, validation ParseValidationFunc, start, end int) HlsParser {
	return HlsParser{
		r:            r,
		validateTS:   validation,
		rootEndpoint: strings.Replace(endpoint, "https:/", "https://", 1),
		uriChan:      make(chan string),
		start:        start,
		end:          end,
	}
}

func (hp HlsParser) Parse() error {
	defer close(hp.uriChan)

	scanner := bufio.NewScanner(hp.r)

	var secondsSum int

	for scanner.Scan() {
		text := scanner.Text()

		if match := timeRegex.FindString(text); match != "" {
			delta, err := hp.timeSum(match)

			if err != nil {
				return err
			}

			secondsSum += delta
		}

		skip, end := hp.TimeSkipper(secondsSum)

		if skip || !hp.validateTS(text) {
			continue
		}

		if end {
			break
		}

		if uri, err := url.Parse(text); err == nil {
			hp.uriChan <- fmt.Sprintf("%s/%s", hp.rootEndpoint, uri.Path[1:])

		} else {
			return err
		}

	}

	hp.uriChan <- ""
	return nil
}

func (hp HlsParser) timeSum(match string) (int, error) {
	switch d, err := hp.textToSeconds(match[:strings.Index(match, ".")]); err {
	case nil:
		return d, nil

	default:
		return -1, err
	}
}

func (hp HlsParser) textToSeconds(sec string) (int, error) {
	d, err := strconv.Atoi(sec)

	if err != nil {
		return 0, err
	}

	return d, nil
}

func (hp HlsParser) TimeSkipper(secondsSum int) (skip, end bool) {
	if secondsSum < hp.start {
		skip = true
	}

	if secondsSum > hp.end && hp.end != 0 {
		end = true
	}

	return
}

func (hp HlsParser) UriChan() chan string {
	return hp.uriChan
}
