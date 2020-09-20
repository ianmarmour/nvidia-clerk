package alert

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

//SendText Sends an SMS notification using Twilio Service.
func SendText(item string, config config.TwilioConfig, client *http.Client) error {
	endpoint := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages", config.ServiceSID)

	msgData := url.Values{}
	msgData.Set("To", config.DestinationNumber)
	msgData.Set("From", config.SourceNumber)
	msgData.Set("Body", fmt.Sprintf("%s Ready for Purchase", item))
	msgDataReader := *strings.NewReader(msgData.Encode())

	req, _ := http.NewRequest("POST", endpoint, &msgDataReader)
	req.SetBasicAuth(config.AccountSID, config.Token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
			return err
		}
	} else {
		fmt.Println(resp.Body)
	}

	return nil
}
