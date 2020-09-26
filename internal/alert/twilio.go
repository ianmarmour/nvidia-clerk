package alert

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

//SendText Sends an SMS notification using Twilio Service.
func SendText(item string, nvidiaURL string, config config.TwilioConfig, client *http.Client) error {
	str := item + " Ready for Purchase: " + "." + nvidiaURL + "."
	api := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages", config.AccountSID)
	data := url.Values{
		"To":   {config.DestinationNumber},
		"From": {config.SourceNumber},
		"Body": {str},
	}
	reader := *strings.NewReader(data.Encode())

	req, err := http.NewRequest("POST", api, &reader)
	if err != nil {
		return err
	}

	req.SetBasicAuth(config.AccountSID, config.Token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	r, err := client.Do(req)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	return nil
}
