package m3u8

type HlsParserInterface interface {
	Parse() error
	UriChan() chan string
}