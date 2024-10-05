package main

import (
	"time"
)

func TimeToSeconds(t string) (int64, error) {
	d, err := time.ParseDuration(t)
	return int64(d.Seconds()), err
}