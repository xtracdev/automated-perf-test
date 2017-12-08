package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"net/http"
)

func TestStartUIMode(t *testing.T) {
	var responseType int = startUiMode()
	assert.Equal(t, http.StatusOK , responseType)
}