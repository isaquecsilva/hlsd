package uribuilder

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(test *testing.T) {
	test.Run("only path",  func(t *testing.T) {
		const (
			playlist = "https://hls.media.net/v1/user/c0r3/vid/aa1bf4ff00/1080_30/playlist.m3u8"
			tsLine =  "/../media/1.ts"
		)

		uri := Build(
			playlist, 
			tsLine,
		)

		expected := "https://" + path.Dir(playlist[8:]) + tsLine
		assert.Equal(t, expected, uri.String())
	})

	test.Run("path and query string", func(t *testing.T) {
		const (
			playlist = "https://hls.media.net/v1/user/c0r3/vid/aa1bf4ff00/1080_30/playlist.m3u8"
			tsLine =  "/../media/1.ts?hash=9ba657581eb180753d7a5f59d3ff4a20c22cbb3c&sig=w9Fhfqb0xOmXNCdNWo6zcbliuIZJEVSbqLreb/clDcuuJo7HgE+8zOvqI6nHZiWaxMV/EEK0MbY4yVRLwAg1jg%%3D%%3D"
		)

		uri := Build(
			playlist, 
			tsLine,
		)

		expected := "https://" + path.Dir(playlist[8:]) + tsLine
		assert.Equal(t, expected, uri.String())
	})

	test.Run("its own host and path", func(t *testing.T) {
		const (
			playlist = "https://hls.media.net/v1/user/c0r3/vid/aa1bf4ff00/1080_30/playlist.m3u8"
			tsLine = "https://hls.media2.net/v2/user/c0r3/vid/aa1bf4ff00/1080_30/1.ts"
		)

		uri := Build(
			playlist,
			tsLine,
		)

		expected := tsLine
		assert.Equal(t, expected, uri.String())
	})

	test.Run("its own host, path and query string", func(t *testing.T) {
		const (
			playlist = "https://hls.media.net/v1/user/c0r3/vid/aa1bf4ff00/1080_30/playlist.m3u8"
			tsLine = "https://hls.media2.net/v2/user/c0r3/vid/aa1bf4ff00/1080_30/1.ts?hash=9ba657581eb180753d7a5f59d3ff4a20c22cbb3c&sig=w9Fhfqb0xOmXNCdNWo6zcbliuIZJEVSbqLreb/clDcuuJo7HgE+8zOvqI6nHZiWaxMV/EEK0MbY4yVRLwAg1jg%%3D%%3D"
		)

		uri := Build(
			playlist,
			tsLine,
		)

		expected := tsLine
		assert.Equal(t, expected, uri.String())
	})	
}