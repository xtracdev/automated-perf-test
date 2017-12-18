package services

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFilePathExist(t *testing.T) {
	path := "C:/Users/a615194/go/src/github.com/xtracdev/automated-perf-test"
	actual := false
	actual = FilePathExist(path)
	expected := true
	assert.Equal (t, expected, actual)
}

func TestFilePathDoesNotExist(t *testing.T) {
	path := "C:/Users/a615194/go/src/github.com/xtrac"
	actual := FilePathExist(path)
	expected := false
	assert.Equal (t, expected, actual)
}
