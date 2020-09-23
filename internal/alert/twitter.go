package alert

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

//SendTweet Sends an Tweet.
func SendTweet(item string, config config.TwitterConfig) error {
	oauth := oauth1.NewConfig(config.ConsumerKey, config.ConsumerSecret)
	token := oauth1.NewToken(config.AccessToken, config.AccessSecret)
	http := oauth.Client(oauth1.NoContext, token)
	twitter := twitter.NewClient(http)

	_, _, err := twitter.Statuses.Update(fmt.Sprintf("%s is now in stock on the NVIDIA Store", item), nil)
	if err != nil {
		return err
	}

	return nil
}
