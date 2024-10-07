package uribuilder

import (
	"net/url"
	"path"
	"strings"

	"github.com/isaquecsilva/hlsd/config"
)

func Build(playlistUri, tsString string) *url.URL {
	uri, err := url.Parse(tsString)

	if err != nil {
		config.Logger().Error(err.Error(), "op", "building-uri")
		return nil
	}

	if uri.Host != "" {
		return uri
	}

	if tsString[0] == '/' {
		tsString = tsString[1:]
	}

	uri, _ = url.ParseRequestURI(strings.Replace(path.Dir(playlistUri), ":/", "://", 1) + "/" + tsString)	
	return uri
}
