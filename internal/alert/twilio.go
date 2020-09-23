package alert

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

//SendText Sends an SMS notification using Twilio Service.
func SendText(item string, config config.TwilioConfig, client *http.Client) error {
	api := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages", config.AccountSID)
	data := url.Values{
		"To":   {config.DestinationNumber},
		"From": {config.SourceNumber},
		"Body": {fmt.Sprintf("%s Ready for Purchase", item)},
	}
	reader := *strings.NewReader(data.Encode())

	req, err := http.NewRequest("POST", api, &reader)
	if err != nil {
		return err
	}

	req.SetBasicAuth(config.AccountSID, config.Token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
