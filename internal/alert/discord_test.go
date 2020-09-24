package alert

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

func TestSendDiscordMessage(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.String() == "http://testurl/webhook/" {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: 503,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
			Header:     make(http.Header),
		}
	})

	cfg := config.DiscordConfig{
		WebhookURL: "http://testurl/webhook/",
	}

	err := SendDiscordMessage("FAKE_SKU_NUMBER", "fakeUrl", cfg, client)
	if err != nil {
		t.Errorf(err.Error())
	}
}
