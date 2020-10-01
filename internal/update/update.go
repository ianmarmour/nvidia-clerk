package update

import (
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/inconshreveable/go-update"
)

// FetchApply Fetches and applys any updates from GitHub releases to this program in place.
func FetchApply(url string, wg *sync.WaitGroup) error {
	defer wg.Done()

	for {
		log.Println("Attempting to fetch updates from github")
		doUpdate(url)
		sleep(60000)
	}
}

func doUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		log.Println("Error applying updates from github.")
	}
	return err
}

func sleep(delay int64) {
	// Force a randomized jitter of up to 5 seconds to avoid looking like a bot.
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(5)

	ns := time.Duration(n) * time.Second
	ds := time.Duration(delay/1000) * time.Second
	time.Sleep(time.Duration(ns + ds))
}
