package alert

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

//SendTweet Sends an Tweet.
func SendTweet(item string, config config.TwitterConfig) error {
	oauthConfig := oauth1.NewConfig(config.ConsumerKey, config.ConsumerSecret)
	token := oauth1.NewToken(config.AccessToken, config.AccessSecret)
	httpClient := oauthConfig.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	tweet, _, err := client.Statuses.Update(item, nil)
	if err != nil {
		fmt.Printf("Error Posting Tweet")
		return err
	}

	fmt.Printf("Posted Tweet\n%v\n", tweet)

	return nil
}
