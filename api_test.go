package main

import (
	"net/http"
	"testing"
)

func TestStartsAPI(t *testing.T) {
	t.Skip()
	leader := Leader()
	go API(leader)
	resp, err := http.Get("http://localhost:9000/api")
	if err != nil {
		t.Fail()
	}

	if resp.StatusCode != 200 {
		t.Fail()
	}
}
