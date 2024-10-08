package m3u8

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/grafov/m3u8"
	"github.com/stretchr/testify/assert"
)

func Test_lines(t *testing.T) {
	parser := NewHlsParser(strings.NewReader("abc\ndef"), "", DefaultParseValidationFunc, 0, 0)

	lines := parser.lines()
	assert.NotNil(t, lines)
	assert.IsType(t, make(chan string), lines)
	for range lines {
	}
}

func Test_isStringTime(test *testing.T) {
	parser := NewHlsParser(bytes.NewBuffer(nil), "", DefaultParseValidationFunc, 0, 0)

	test.Run("true // is time", func(t *testing.T) {
		timestamp, isTime := parser.isStringTime("#EXTINF:3.003000,")
		assert.True(t, isTime)
		assert.Equal(t, "3.003000", timestamp)
	})

	test.Run("false // is not time", func(t *testing.T) {
		timestamp, isTime := parser.isStringTime("/v1/vod/100.ts")
		assert.False(t, isTime)
		assert.Empty(t, timestamp)
	})
}

func Test_textToSeconds(test *testing.T) {
	parser := NewHlsParser(bytes.NewBuffer(nil), "", DefaultParseValidationFunc, 0, 0)

	test.Run("error - not a valid timestamp expression", func(t *testing.T) {
		timestamp := parser.textToSeconds("abc")
		assert.Zero(t, timestamp)
	})

	test.Run("ok - valid timestamp", func(t *testing.T) {
		timestamp := parser.textToSeconds("15.00000")
		assert.Equal(t, int64(15), timestamp)
	})
}

func Test_timeSkipper(test *testing.T) {
	parser := NewHlsParser(bytes.NewBuffer(nil), "", DefaultParseValidationFunc, 0, 0)

	test.Run("no skip, no end", func(t *testing.T) {
		secondsSum := int64(60)
		skip, end := parser.timeSkipper(secondsSum)
		assert.False(t, skip)
		assert.False(t, end)
	})

	test.Run("skip ten seconds, no end", func(t *testing.T) {
		parser.start = 10

		secondsSum := int64(0)
		skip, end := parser.timeSkipper(secondsSum)
		assert.True(t, skip)
		assert.False(t, end)
	})

	test.Run("no skip, end 20", func(t *testing.T) {
		parser.start = 0
		parser.end = 20

		secondsSum := int64(20)

		skip, end := parser.timeSkipper(secondsSum)
		assert.False(t, skip)
		assert.True(t, end)
	})

	test.Run("skip 20, end 30", func(t *testing.T) {
		parser.start = 20
		parser.end = 30

		tests := []struct {
			seconds   int64
			skip, end bool
		}{
			{
				seconds: 0,
				skip:    true,
				end:     false,
			},
			{
				seconds: 10,
				skip:    true,
				end:     false,
			},
			{
				seconds: 20,
				skip:    false,
				end:     false,
			},
			{
				seconds: 30,
				skip:    false,
				end:     true,
			},
		}

		for _, test := range tests {
			skip, end := parser.timeSkipper(test.seconds)
			assert.Equal(t, test.skip, skip)
			assert.Equal(t, test.end, end)
		}
	})

}

func TestParse(t *testing.T) {
	r := m3u8FileGenerator(20)

	parser := NewHlsParser(r, "https://root.endpoint.net/chunks/playlist.m3u8", DefaultParseValidationFunc, 25, 0)

	go parser.Parse()

	links := []string{}

	for uri := range parser.UriChan(){
		links = append(links, uri)
	}

	assert.Len(t, links, 18)

	for _, uri := range links {
		assert.Contains(t, uri, ".ts")
	}
}

func m3u8FileGenerator(links uint) io.Reader {
	p, err := m3u8.NewMediaPlaylist(links, links)

	if err != nil {
		log.Fatal(err)
	}

	for i := range 20 {
		p.Append(fmt.Sprintf("/%d.ts", i), 10, "")
	}

	return p.Encode()
}
