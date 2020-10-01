package rest

import (
	"encoding/json"
	"net/http"
)

type Config struct {
	Sections []Section `json:"sections"`
}

//JSON Returns the Marshalled Version of the Response
func (r *Config) JSON() ([]byte, error) {
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

type Section struct {
	ID         int64       `json:"id"`
	Name       string      `json:"name"`
	Components []Component `json:"components"`
}

type Component struct {
	ID    int64       `json:"id"`
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type Status struct {
	Healthy bool `json:"status"`
}

//JSON Returns the Marshalled Version of the Response
func (r *Status) JSON() ([]byte, error) {
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func StatusResponse(w http.ResponseWriter, r *http.Request) {
	resBody := Status{
		Healthy: true,
	}

	json, _ := resBody.JSON()
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func ConfigResponseMain(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Declare a new Person struct.
		var sect Section

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&sect)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}
}

func ConfigResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		resBody := Config{
			Sections: []Section{
				{
					ID:   15,
					Name: "Main",
					Components: []Component{
						{
							ID:    1,
							Name:  "Region",
							Type:  "select",
							Value: []string{"GBR", "USA"},
						},
						{
							ID:    1,
							Name:  "Model",
							Type:  "select",
							Value: []string{"3080", "3090"},
						},
					},
				},
				{
					ID:   15,
					Name: "Discord",
					Components: []Component{
						{
							ID:   1,
							Name: "Webhook URL",
							Type: "input",
						},
					},
				},
				{
					ID:   15,
					Name: "Twilio",
					Components: []Component{
						{
							ID:   1,
							Name: "Source",
							Type: "input",
						},
						{
							ID:   1,
							Name: "Destination",
							Type: "input",
						},
						{
							ID:   1,
							Name: "API Key",
							Type: "input",
						},
						{
							ID:   1,
							Name: "Account SID",
							Type: "input",
						},
					},
				},
				{
					ID:   15,
					Name: "Telegram",
					Components: []Component{
						{
							ID:   1,
							Name: "API Key",
							Type: "input",
						},
						{
							ID:   1,
							Name: "Chat ID",
							Type: "input",
						},
					},
				},
			},
		}

		json, _ := resBody.JSON()
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}

	if r.Method == http.MethodPost {

	}
}
