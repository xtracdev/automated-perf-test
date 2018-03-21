package services

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilePathExist(t *testing.T) {
	path := os.Getenv("GOPATH") + "/src/github.com/xtracdev/automated-perf-test/uiServices/test/"
	actual := false
	fmt.Println(path)
	actual = FilePathExist(path)
	expected := true
	assert.Equal(t, expected, actual)
}

func TestFilePathDoesNotExist(t *testing.T) {
	path := "((((((("
	actual := FilePathExist(path)
	expected := false
	assert.Equal(t, expected, actual)
}
