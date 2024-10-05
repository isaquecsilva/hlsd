package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeToSeconds(test *testing.T) {

	test.Run("success parsing", func(t *testing.T) {
		// success	
		seconds, err := TimeToSeconds("01h20m30s")
		assert.Nil(t, err)
		assert.Equal(t, int64(4830), seconds)
	})

	test.Run("error invalid time string", func(t *testing.T) {
		seconds, err := TimeToSeconds("01:20:30")
		assert.Zero(t, seconds)
		assert.NotNil(t, err)
	})

}