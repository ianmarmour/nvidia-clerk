package alert

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

func TestSendText(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		if req.URL.String() == "https://api.twilio.com/2010-04-01/Accounts/1/Messages" {
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

	cfg := config.TwilioConfig{
		AccountSID:        "1",
		Token:             "fake",
		SourceNumber:      "fake",
		DestinationNumber: "fake",
	}

	err := SendText("FAKE_SKU_NUMBER", cfg, client)
	if err != nil {
		t.Errorf(err.Error())
	}
}
