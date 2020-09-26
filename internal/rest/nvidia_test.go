package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestGetSessionToken(t *testing.T) {
	token := SessionToken{"12345"}
	json, err := json.Marshal(token)

	client := NewTestClient(func(req *http.Request) *http.Response {
		if strings.Contains(req.URL.String(), "https://store.nvidia.com/store/nvidia/SessionToken?format=json") {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader(json)),
				Header:     make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: 503,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
			Header:     make(http.Header),
		}
	})

	sessionToken, err := GetSessionToken(client)
	if err != nil {
		t.Errorf(err.Error())
	}

	if sessionToken.Value != "12345" {
		t.Errorf(err.Error())
	}
}
