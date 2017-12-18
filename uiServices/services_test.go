package services

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
)



func TestFilePathExist(t *testing.T) {
	path := os.Getenv("GOPATH")
	actual := false
	actual = FilePathExist(path)
	expected := true
	os.Getenv("GOPATH")
	assert.Equal (t, expected, actual)
}

func TestFilePathDoesNotExist(t *testing.T) {
	path := "C:/Users/a615194/go/src/github.com/xtrac"
	actual := FilePathExist(path)
	expected := false
	assert.Equal (t, expected, actual)
}

//
//func TestFileSave(t *testing.T) {
//
//	path :=  "C:/Users/a615194/go/src/github.com/xtracdev/automated-perf-test/config/ConfigTest.xml"
//	actual := FilePathExist(path)
//	expected := false
//	assert.Equal (t, expected, actual)
//}
