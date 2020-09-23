package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// resetEnv Simple function to reset environment to existing environment.
func resetEnv(env []string) func() {
	return func() {
		for _, e := range env {
			pair := strings.SplitN(e, "=", 2)
			os.Setenv(pair[0], pair[1])
		}
	}
}

func envSMS() func() {
	vars := []string{"TWILIO_ACCOUNT_SID=1", "TWILIO_TOKEN=2", "TWILIO_SOURCE_NUMBER=3", "TWILIO_DESTINATION_NUMBER=4"}

	return func() {
		for _, e := range vars {
			pair := strings.SplitN(e, "=", 2)
			os.Setenv(pair[0], pair[1])
		}
	}
}

func envDiscord() func() {
	vars := []string{"DISCORD_WEBHOOK_URL=1"}

	return func() {
		for _, e := range vars {
			pair := strings.SplitN(e, "=", 2)
			os.Setenv(pair[0], pair[1])
		}
	}
}

func envTwitter() func() {
	vars := []string{"TWITTER_CONSUMER_KEY=1", "TWITTER_CONSUMER_SECRET=2", "TWITTER_ACCESS_TOKEN=3", "TWITTER_ACCESS_SECRET=4"}

	return func() {
		for _, e := range vars {
			pair := strings.SplitN(e, "=", 2)
			os.Setenv(pair[0], pair[1])
		}
	}
}

func envTelegram() func() {
	vars := []string{"TELEGRAM_API_KEY=1", "TELEGRAM_CHAT_ID=2"}

	return func() {
		for _, e := range vars {
			pair := strings.SplitN(e, "=", 2)
			os.Setenv(pair[0], pair[1])
		}
	}
}

func TestGet(t *testing.T) {
	tests := map[string]struct {
		region      string
		delay       int64
		sms         bool
		discord     bool
		twitter     bool
		telegram    bool
		environment func()
		expected    *Config
	}{
		"default": {
			region:      "USA",
			delay:       500,
			sms:         false,
			discord:     false,
			twitter:     false,
			telegram:    false,
			environment: func() {},
			expected: &Config{
				Locale:   "en_us",
				Currency: "USD",
				Delay:    500,
				SKU:      strPtr("5438481700"),
				TestSKU:  strPtr("5379432500"),
			},
		},
		"with sms": {
			region:      "USA",
			sms:         true,
			discord:     false,
			twitter:     false,
			telegram:    false,
			environment: envSMS(),
			expected: &Config{
				Locale:   "en_us",
				Currency: "USD",
				Delay:    0,
				SKU:      strPtr("5438481700"),
				TestSKU:  strPtr("5379432500"),
				TwilioConfig: TwilioConfig{
					AccountSID:        "1",
					Token:             "2",
					SourceNumber:      "3",
					DestinationNumber: "4",
				},
			},
		},
		"with discord": {
			region:      "USA",
			sms:         false,
			discord:     true,
			twitter:     false,
			telegram:    false,
			environment: envDiscord(),
			expected: &Config{
				Locale:   "en_us",
				Currency: "USD",
				Delay:    0,
				SKU:      strPtr("5438481700"),
				TestSKU:  strPtr("5379432500"),
				DiscordConfig: DiscordConfig{
					WebhookURL: "1",
				},
			},
		},
		"with twitter": {
			region:      "USA",
			sms:         false,
			discord:     false,
			twitter:     true,
			telegram:    false,
			environment: envTwitter(),
			expected: &Config{
				Locale:   "en_us",
				Currency: "USD",
				Delay:    0,
				SKU:      strPtr("5438481700"),
				TestSKU:  strPtr("5379432500"),
				TwitterConfig: TwitterConfig{
					ConsumerKey:    "1",
					ConsumerSecret: "2",
					AccessToken:    "3",
					AccessSecret:   "4",
				},
			},
		},
		"with telegram": {
			region:      "USA",
			sms:         false,
			discord:     false,
			twitter:     false,
			telegram:    true,
			environment: envTelegram(),
			expected: &Config{
				Locale:   "en_us",
				Currency: "USD",
				Delay:    0,
				SKU:      strPtr("5438481700"),
				TestSKU:  strPtr("5379432500"),
				TelegramConfig: TelegramConfig{
					APIKey: "1",
					ChatID: "2",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			defer resetEnv(os.Environ())

			test.environment()

			result, err := Get(test.region, "3080", test.delay, test.sms, test.discord, test.twitter, test.telegram)
			if err != nil {
				t.Errorf(err.Error())
			}
			assert.Equal(t, test.expected, result)
		})
	}
}
