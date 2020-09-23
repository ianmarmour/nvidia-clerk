package alert

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

func TestSendTelegramMessage(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.String() == "https://api.telegram.org/bot1/sendMessage" {
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

	cfg := config.TelegramConfig{
		APIKey: "1",
		ChatID: "1",
	}

	err := SendTelegramMessage("FAKE_SKU_NUMBER", cfg, client)
	if err != nil {
		t.Errorf(err.Error())
	}
}
