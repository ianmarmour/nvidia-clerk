package config

import (
	"fmt"
	"os"
)

func strPtr(in string) *string {
	i := in
	return &i
}

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

type Model struct {
	SKU *string
}

type RegionalConfig struct {
	Models       map[string]Model
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
	Locale       string
	NvidiaLocale string
	Currency     string
	Delay        int64

	SKU            *string
	TestSKU        *string
	TwilioConfig   *TwilioConfig
	TwitterConfig  *TwitterConfig
	DiscordConfig  *DiscordConfig
	TelegramConfig *TelegramConfig
}

// Hardcoded SKU to locale/currency mappings to avoid user pain of having to lookup and enter these.
var regionalConfig = map[string]RegionalConfig{
	"AUT": {
		Models: map[string]Model{
			"3070": {},
			"3080": {
				SKU: strPtr("5440853700"),
			},
			"3090": {
				SKU: strPtr("5444941400"),
			},
		},
		Locale:       "de_de",
		NvidiaLocale: "de-de",
		Currency:     "EUR",
	},
	"BEL": {
		Models: map[string]Model{
			"3070": {
				SKU: nil,
			},
			"3080": {
				SKU: strPtr("5438795700"),
			},
			"3090": {
				SKU: strPtr("5438795600"),
			},
		},
		Locale:       "fr_fr",
		NvidiaLocale: "fr-fr",
		Currency:     "EUR",
	},
	"CAN": {
		Models: map[string]Model{
			"3070": {},
			"3080": {
				SKU: strPtr("5438481700"),
			},
			"3090": {
				SKU: strPtr("5438481600"),
			},
		},
		Locale:       "en_us",
		NvidiaLocale: "en-ca",
		Currency:     "CAN",
	},
	"CZE": {
		Models: map[string]Model{
			"3070": {},
			"3080": {
				SKU: strPtr("5438793800"),
			},
			"3090": {
				SKU: strPtr("5438793600"),
			},
		},
		Locale:       "en_gb",
		NvidiaLocale: "en-gb",
		Currency:     "EUR",
	},
	"DNK": {
		Models: map[string]Model{
			"3070": {},
			"3080": {
				SKU: strPtr("5438793300"),
			},
			"3090": {
				SKU: strPtr("5438793500"),
			},
		},
		Locale:       "en_gb",
		NvidiaLocale: "en-gb",
		Currency:     "EUR",
	},
	"FIN": {
		Models: map[string]Model{
			"3070": {},
			"3080": {
				SKU: strPtr("5438793300"),
			},
			"3090": {
				SKU: strPtr("5438793500"),
			},
		},
		Locale:       "en_gb",
		NvidiaLocale: "en-gb",
		Currency:     "EUR",
	},
	"FRA": {
		Models: map[string]Model{
			"3070": {},
			"3080": {
				SKU: strPtr("5438795200"),
			},
			"3090": {
				SKU: strPtr("5438761500"),
			},
		},
		Locale:       "fr_fr",
		NvidiaLocale: "fr-fr",
		Currency:     "EUR",
	},
	"DEU": {
		Models: map[string]Model{
			"3070": {},
			"3080": {
				SKU: strPtr("5438792300"),
			},
			"3090": {
				SKU: strPtr("5438761400"),
			},
		},
		Locale:       "de_de",
		NvidiaLocale: "de-de",
		Currency:     "EUR",
	},
	"USA": {
		Models: map[string]Model{
			"2060": {
				SKU: strPtr("5379432500"),
			},
			"3070": {},
			"3080": {
				SKU: strPtr("5438481700"),
			},
			"3090": {
				SKU: strPtr("5438481600"),
			},
		},
		Locale:       "en_us",
		NvidiaLocale: "en-us",
		Currency:     "USD",
		TestSKU:      "5379432500",
	},
	"GBR": {
		Models: map[string]Model{
			"3070": {},
			"3080": {
				SKU: strPtr("5438792800"),
			},
			"3090": {
				SKU: strPtr("5438792700"),
			},
		},
		Locale:       "en_gb",
		NvidiaLocale: "en-gb",
		Currency:     "GBP",
	},
	"IRL": {
		Models: map[string]Model{
			"3070": {},
			"3080": {
				SKU: strPtr("5438792800"),
			},
			"3090": {
				SKU: strPtr("5438792700"),
			},
		},
		Locale:       "en_gb",
		NvidiaLocale: "en-gb",
		Currency:     "GBP",
	},
	"ITA": {
		Models: map[string]Model{
			"3070": {},
			"3080": {
				SKU: strPtr("5438796200"),
			},
			"3090": {
				SKU: strPtr("5438796100"),
			},
		},
		Locale:       "it_it",
		NvidiaLocale: "it-it",
		Currency:     "EUR",
	},
	"SWE": {
		Models: map[string]Model{
			"3070": {},
			"3080": {
				SKU: strPtr("5438798100"),
			},
			"3090": {
				SKU: strPtr("5438761600"),
			},
		},
		Locale:       "sv_SE",
		NvidiaLocale: "sv-se",
		Currency:     "SEK",
	},
	"LUX": {
		Models: map[string]Model{
			"3070": {},
			"3080": {
				SKU: strPtr("5438795700"),
			},
			"3090": {
				SKU: strPtr("5438795600"),
			},
		},
		Locale:       "fr_fr",
		NvidiaLocale: "fr-fr",
		Currency:     "EUR",
	},
	"POL": {
		Models: map[string]Model{
			"3070": {},
			"3080": {
				SKU: strPtr("5438797700"),
			},
			"3090": {
				SKU: strPtr("5438797600"),
			},
		},
		Locale:       "pl_pl",
		NvidiaLocale: "pl-pl",
		Currency:     "PLN",
	},
	"PRT": {
		Models: map[string]Model{
			"3070": {},
			"3080": {
				SKU: strPtr("5438794300"),
			},
			"3090": {},
		},
		Locale:       "en_gb",
		NvidiaLocale: "en-gb",
		Currency:     "EUR",
	},
	"ESP": {
		Models: map[string]Model{
			"3070": {},
			"3080": {
				SKU: strPtr("5438794800"),
			},
			"3090": {
				SKU: strPtr("5438794700"),
			},
		},
		Locale:       "es_es",
		NvidiaLocale: "es-es",
		Currency:     "EUR",
	},
	"NOR": {
		Models: map[string]Model{
			"3070": {},
			"3080": {
				SKU: strPtr("5438797200"),
			},
			"3090": {
				SKU: strPtr("5438797100"),
			},
		},
		Locale:       "no_NO",
		NvidiaLocale: "no-NO",
		Currency:     "NOK",
	},
	"NLD": {
		Models: map[string]Model{
			"3070": {},
			"3080": {
				SKU: strPtr("5438796700"),
			},
			"3090": {
				SKU: strPtr("5438796600"),
			},
		},
		Locale:       "nl_nl",
		NvidiaLocale: "nl-nl",
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
		return nil, &ConfigError{"Twilio", "TWILIO_ACCOUNT_SID"}
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
func Get(region string, model string, delay int64, sms bool, discord bool, twitter bool, telegram bool) (*Config, error) {
	if regionConfig, ok := regionalConfig[region]; ok {
		configuration := Config{}

		configuration.SKU = regionConfig.Models[model].SKU
		configuration.TestSKU = &regionConfig.TestSKU
		configuration.Delay = delay
		configuration.Locale = regionConfig.Locale
		configuration.NvidiaLocale = regionConfig.NvidiaLocale
		configuration.Currency = regionConfig.Currency

		if sms == true {
			cfg, err := getTwilio()
			if err != nil {
				return nil, err
			}
			configuration.TwilioConfig = cfg
		}

		if discord == true {
			cfg, err := getDiscord()
			if err != nil {
				return nil, err
			}
			configuration.DiscordConfig = cfg
		}

		if twitter == true {
			cfg, err := getTwitter()
			if err != nil {
				return nil, err
			}
			configuration.TwitterConfig = cfg
		}

		if telegram == true {
			cfg, err := getTelegram()
			if err != nil {
				return nil, err
			}
			configuration.TelegramConfig = cfg
		}

		return &configuration, nil
	}

	return nil, &RegionError{region}
}
