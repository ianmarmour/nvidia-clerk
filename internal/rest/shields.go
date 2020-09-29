package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

//ShieldsEndpointResponse Represents a valid endpoint response for shields.io
type shieldsResponse struct {
	Version   int    `json:"schemaVersion"`
	NamedLogo string `json:"namedLogo"`
	Label     string `json:"label"`
	Message   string `json:"message"`
	Color     string `json:"color"`
}

//JSON Returns the Marshalled Version of the Response
func (r *shieldsResponse) JSON() ([]byte, error) {
	payload, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

var client *http.Client

// init Initializes required variables.
func init() {
	client = &http.Client{Timeout: 10 * time.Second}
}

func newShieldsResponse() shieldsResponse {
	return shieldsResponse{
		Version:   1,
		Label:     "nvidia",
		NamedLogo: "nvidia",
		Color:     "green",
	}
}

func getShieldsResponse() []byte {
	res := newShieldsResponse()
	_, err := GetSessionToken(client)
	if err != nil {
		res.Message = "offline"
	}

	res.Message = "online"

	json, err := res.JSON()
	if err != nil {
		log.Println(err)
		return nil
	}

	return json
}

func endpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, string(getShieldsResponse()))
}

//StartShieldsAPIServer Starts up a shields API server
func StartShieldsAPIServer(config config.ShieldsConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/endpoint", endpoint)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), router))
}
