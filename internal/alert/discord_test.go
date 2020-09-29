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

	productMsg := DiscordProductMessage{}
	productMsg.Set("test", "test")

	err := SendDiscordMessage(&productMsg, cfg, client)
	if err != nil {
		t.Errorf(err.Error())
	}

	apiMsg := DiscordAPIMessage{}
	apiMsg.Set("test", "test")

	err = SendDiscordMessage(&apiMsg, cfg, client)
	if err != nil {
		t.Errorf(err.Error())
	}
}
