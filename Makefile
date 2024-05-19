build:
	go build -ldflags='-s -w' -o bin/twitch.exe .
run:
	@if [ -f stream.ts ]; then rm stream.ts; fi
	go run . -m3u8 $(uri) -start 1200 -end 1800 -queue 5

help:
	go run . -help