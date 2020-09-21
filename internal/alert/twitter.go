package alert

import (
	"fmt"
	"log"

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

	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}
	user, _, _ := client.Accounts.VerifyCredentials(verifyParams)
	fmt.Printf("User's ACCOUNT:\n%+v\n", user)

	_, _, err := client.Statuses.Update("A Test Tweet from a new Bot I'm building!", nil)
	if err != nil {
		log.Println("Error attempting to Tweet")
		return err
	}

	return nil
}
