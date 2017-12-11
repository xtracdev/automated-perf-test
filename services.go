package main

import (
"encoding/json"
"log"
"net/http"
)

type test_struct struct {
Test string
}

func test(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t test_struct
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	log.Println(t.Test)
}
