package main

import (
	"net/url"
	"path"
	"strings"
	"sync"
	"testing"
)

func FuzzParse(f *testing.F) {
	var seeds []string = []string {
		"https://some-fake-hosting/cdn/v1/app/11-video-id/i1.ts",
		"https://fake.cdn.net/##$!qq_eiir482fdd/efev/aa/index-00.ts",
		"/vplayer/cc/enus/100_index.ts",
		"/hls/stream-11sso-11994/v.1.ts",
		"/video-1.ts",
		"/50.ts",
		"/--308.ts?auth=fake-key-auth&state=permission&quality=1080p&fps=60",
		"/stream/6c657347-8eff-473d-97af-3f17266c418c/720p60/v77.ts",
	}

	// settings seeds for fuzzing test
	for _, seed := range seeds {
		f.Add(seed)
	}

	var wg sync.WaitGroup

	const defaultEndpoint = "https://fake-cdn-hosting.icnet/file.m3u8"

	var link string

	f.Fuzz(func(t *testing.T, in string) {
		parser := NewHlsParser(strings.NewReader(in), path.Dir(defaultEndpoint), DefaultParseValidationFunc, 0, 0)
		wg.Add(1)

		go func() {
			link = <- parser.UriChan()
			<- parser.UriChan()
			wg.Done()
		}()		

		err := parser.Parse()
		wg.Wait()

		if err != nil {
			t.Error(err)
		}

		uri, _ := url.Parse(defaultEndpoint)
		inPath, _ := url.Parse(in)
		expected := "https://" + uri.Host + inPath.Path

		if expected != link {
			t.Errorf("expected = %v, actual = %v", expected, link)
		}
	})
}