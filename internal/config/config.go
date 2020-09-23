package config

import (
	"fmt"
	"os"
)

type RegionError struct {
	Code string
}

func (w *RegionError) Error() string {
	return fmt.Sprintf("%s: Region unsupported", w.Code)
}

type ConfigError struct {
	Type string
	Name string
}

func (w *ConfigError) Error() string {
	return fmt.Sprintf("%s: %v Environment variable not found", w.Type, w.Name)
}

type RegionalConfig struct {
	SKU          string
	Locale       string
	NvidiaLocale string
	Currency     string
	TestSKU      string
}

type TwitterConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

type TwilioConfig struct {
	AccountSID        string
	Token             string
	SourceNumber      string
	DestinationNumber string
}

type DiscordConfig struct {
	WebhookURL string
}

type TelegramConfig struct {
	APIKey string
	ChatID string
}

type Config struct {
	Locale   string
	Currency string
	Delay    int64

	SKU            string
	TestSKU        string
	TwilioConfig   TwilioConfig
	TwitterConfig  TwitterConfig
	DiscordConfig  DiscordConfig
	TelegramConfig TelegramConfig
}

// Hardcoded SKU to locale/currency mappings to avoid user pain of having to lookup and enter these.
var regionalConfig = map[string]RegionalConfig{
	"AUT": {
		SKU:          "5440853700",
		Locale:       "de_de",
		NvidiaLocale: "de_de",
		Currency:     "EUR",
	},
	"BEL": {
		SKU:          "5438795700",
		Locale:       "fr_fr",
		NvidiaLocale: "fr_fr",
		Currency:     "EUR",
	},
	"CAN": {
		SKU:          "5438481700",
		Locale:       "en_us",
		NvidiaLocale: "en_ca",
		Currency:     "CAN",
	},
	"CZE": {
		SKU:          "5438793800",
		Locale:       "en_gb",
		NvidiaLocale: "en_gb",
		Currency:     "EUR",
	},
	"DNK": {
		SKU:          "5438793300",
		Locale:       "en_gb",
		NvidiaLocale: "en_gb",
		Currency:     "EUR",
	},
	"FIN": {
		SKU:          "5438793300",
		Locale:       "en_gb",
		NvidiaLocale: "en_gb",
		Currency:     "EUR",
	},
	"FRA": {
		SKU:          "5438795200",
		Locale:       "fr_fr",
		NvidiaLocale: "fr_fr",
		Currency:     "EUR",
	},
	"DEU": {
		SKU:          "5438792300",
		Locale:       "de_de",
		NvidiaLocale: "de_de",
		Currency:     "EUR",
	},
	"USA": {
		SKU:          "5438481700",
		Locale:       "en_us",
		NvidiaLocale: "en_us",
		Currency:     "USD",
		TestSKU:      "5379432500",
	},
	"GBR": {
		SKU:          "5438792800",
		Locale:       "en_gb",
		NvidiaLocale: "en_gb",
		Currency:     "GBP",
	},
	"IRL": {
		SKU:          "5438792800",
		Locale:       "en_gb",
		NvidiaLocale: "en_gb",
		Currency:     "GBP",
	},
	"ITA": {
		SKU:          "5438796200",
		Locale:       "it_it",
		NvidiaLocale: "it_it",
		Currency:     "EUR",
	},
	"SWE": {
		SKU:          "5438798100",
		Locale:       "sv_SE",
		NvidiaLocale: "sv_se",
		Currency:     "SEK",
	},
	"LUX": {
		SKU:          "5438795700",
		Locale:       "fr_fr",
		NvidiaLocale: "fr_fr",
		Currency:     "EUR",
	},
	"POL": {
		SKU:          "5438797700",
		Locale:       "pl_pl",
		NvidiaLocale: "pl_pl",
		Currency:     "PLN",
	},
	"PRT": {
		SKU:          "5438794300",
		Locale:       "en_gb",
		NvidiaLocale: "en_gb",
		Currency:     "EUR",
	},
	"ESP": {
		SKU:          "5438794800",
		Locale:       "es_es",
		NvidiaLocale: "es_es",
		Currency:     "EUR",
	},
	"NOR": {
		SKU:          "5438797200",
		Locale:       "no_NO",
		NvidiaLocale: "no_NO",
		Currency:     "NOK",
	},
	"NLD": {
		SKU:          "5438796700",
		Locale:       "nl_nl",
		NvidiaLocale: "nl_nl",
		Currency:     "EUR",
	},
}

//getTwitter Generates TwitterConfiguration for application from environmental variables.
func getTwitter() (*TwitterConfig, error) {
	c := TwitterConfig{}

	key, keyOk := os.LookupEnv("TWITTER_CONSUMER_KEY")
	if keyOk == false {
		return nil, &ConfigError{"Twitter", "TWITTER_CONSUMER_KEY"}
	}
	c.ConsumerKey = key

	cs, csOk := os.LookupEnv("TWITTER_CONSUMER_SECRET")
	if csOk == false {
		return nil, &ConfigError{"Twitter", "TWITTER_CONSUMER_SECRET"}
	}
	c.ConsumerSecret = cs

	at, atOk := os.LookupEnv("TWITTER_ACCESS_TOKEN")
	if atOk == false {
		return nil, &ConfigError{"Twitter", "TWITTER_ACCESS_TOKEN"}
	}

	c.AccessToken = at

	as, asOk := os.LookupEnv("TWITTER_ACCESS_SECRET")
	if asOk == false {
		return nil, &ConfigError{"Twitter", "TWITTER_ACCESS_SECRET"}
	}

	c.AccessSecret = as

	return &c, nil
}

//getTwilio Generates TwilioConfiguration for application from environmental variables.
func getTwilio() (*TwilioConfig, error) {
	c := TwilioConfig{}

	id, idOk := os.LookupEnv("TWILIO_ACCOUNT_SID")
	if idOk == false {
		return nil, &ConfigError{"Twilio", "TWLIO_ACCOUNT_SID"}
	}
	c.AccountSID = id

	t, tOk := os.LookupEnv("TWILIO_TOKEN")
	if tOk == false {
		return nil, &ConfigError{"Twilio", "TWILIO_TOKEN"}
	}
	c.Token = t

	sn, snOk := os.LookupEnv("TWILIO_SOURCE_NUMBER")
	if snOk == false {
		return nil, &ConfigError{"Twilio", "TWILIO_SOURCE_NUMBER"}
	}
	c.SourceNumber = sn

	dn, dnOk := os.LookupEnv("TWILIO_DESTINATION_NUMBER")
	if dnOk == false {
		return nil, &ConfigError{"Twilio", "TWILIO_DESTINATION_NUMBER"}
	}
	c.DestinationNumber = dn

	return &c, nil
}

//getDiscord Generates DiscordConfiguration for application from environmental variables.
func getDiscord() (*DiscordConfig, error) {
	c := DiscordConfig{}

	u, uOk := os.LookupEnv("DISCORD_WEBHOOK_URL")
	if uOk == false {
		return nil, &ConfigError{"Discord", "DISCORD_WEBHOOK_URL"}
	}
	c.WebhookURL = u

	return &c, nil
}

//getTelegram Generates TelegramConfiguration for application from environmental variables.
func getTelegram() (*TelegramConfig, error) {
	c := TelegramConfig{}

	a, aOk := os.LookupEnv("TELEGRAM_API_KEY")
	if aOk == false {
		return nil, &ConfigError{"Telegram", "TELEGRAM_API_KEY"}
	}
	c.APIKey = a

	id, idOk := os.LookupEnv("TELEGRAM_CHAT_ID")
	if idOk == false {
		return nil, &ConfigError{"Telegram", "TELEGRAM_CHAT_ID"}
	}

	c.ChatID = id

	return &c, nil
}

//Get Generates Configuration for application from environmental variables.
func Get(region string, delay int64, sms bool, discord bool, twitter bool, telegram bool) (*Config, error) {
	if regionConfig, ok := regionalConfig[region]; ok {
		configuration := Config{}

		configuration.SKU = regionConfig.SKU
		configuration.Delay = delay
		configuration.TestSKU = regionConfig.TestSKU
		configuration.Locale = regionConfig.Locale
		configuration.Currency = regionConfig.Currency

		if sms == true {
			cfg, err := getTwilio()
			if err != nil {
				return nil, err
			}
			configuration.TwilioConfig = *cfg
		}

		if discord == true {
			cfg, err := getDiscord()
			if err != nil {
				return nil, err
			}
			configuration.DiscordConfig = *cfg
		}

		if twitter == true {
			cfg, err := getTwitter()
			if err != nil {
				return nil, err
			}
			configuration.TwitterConfig = *cfg
		}

		if telegram == true {
			cfg, err := getTelegram()
			if err != nil {
				return nil, err
			}
			configuration.TelegramConfig = *cfg
		}

		return &configuration, nil
	}

	return nil, &RegionError{region}
}
