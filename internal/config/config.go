package config

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

type RegionError struct {
	Code string
}

func (w *RegionError) Error() string {
	return fmt.Sprintf("%s: region unsupported", w.Code)
}

type ModelError struct {
	Code string
}

func (w *ModelError) Error() string {
	return fmt.Sprintf("%s: model unsupported", w.Code)
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

type ToastConfig struct {
	OS string
}

type RegionalConfig struct {
	Models       map[string]Model
	Locale       string
	NvidiaLocale string
	Currency     string
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

type ShieldsConfig struct {
	Port string
}

type Config struct {
	Locale       string
	NvidiaLocale string
	Currency     string
	Delay        int64

	SKU            *string
	TwilioConfig   *TwilioConfig
	TwitterConfig  *TwitterConfig
	DiscordConfig  *DiscordConfig
	TelegramConfig *TelegramConfig
	ToastConfig    *ToastConfig
	ShieldsConfig  *ShieldsConfig
}

// Hardcoded SKU to locale/currency mappings to avoid user pain of having to lookup and enter these.
var RegionalConfigs = map[string]RegionalConfig{
	"AUT": {
		Models: map[string]Model{
			"2060": {
				SKU: strPtr("5394902900"),
			},
			"2070": {
				SKU: strPtr("5394901600"),
			},
			"2080": {
				SKU: strPtr("5335703700"),
			},
			"2080TI": {
				SKU: strPtr("5218984600"),
			},
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
			"2060": {
				SKU: strPtr("5394902700"),
			},
			"2070": {
				SKU: strPtr("5336534300"),
			},
			"2080": {
				SKU: strPtr("5336531500"),
			},
			"2080TI": {
				SKU: strPtr("5218987100"),
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
			"2060": {
				SKU: strPtr("5379432500"),
			},
			"2070": {
				SKU: strPtr("5379432400"),
			},
			"2080TI": {
				SKU: strPtr("5218984100"),
			},
			"3080": {
				SKU: strPtr("5438481700"),
			},
			"3090": {
				SKU: strPtr("5438481600"),
			},
		},
		Locale:       "en_us",
		NvidiaLocale: "en-us",
		Currency:     "CAD",
	},
	"CZE": {
		Models: map[string]Model{
			"2060": {
				SKU: strPtr("5394902800"),
			},
			"2070": {
				SKU: strPtr("5394901500"),
			},
			"2080": {
				SKU: strPtr("5336531900"),
			},
			"2080TI": {
				SKU: strPtr("5218613300"),
			},
			"3080": {
				SKU: strPtr("5438793800"),
			},
			"3090": {
				SKU: strPtr("5438793600"),
			},
		},
		Locale:       "en_gb",
		NvidiaLocale: "en-gb",
		Currency:     "CZK",
	},
	"DNK": {
		Models: map[string]Model{
			"2060": {
				SKU: strPtr("5394903100"),
			},
			"2070": {
				SKU: strPtr("5394901800"),
			},
			"2080": {
				SKU: strPtr("5336531800"),
			},
			"2080TI": {
				SKU: strPtr("5218988600"),
			},
			"3080": {
				SKU: strPtr("5438793300"),
			},
			"3090": {
				SKU: strPtr("5438793500"),
			},
		},
		Locale:       "en_gb",
		NvidiaLocale: "en-gb",
		Currency:     "DKK",
	},
	"FIN": {
		Models: map[string]Model{
			"2060": {
				SKU: strPtr("5394903100"),
			},
			"2070": {
				SKU: strPtr("5394901800"),
			},
			"2080": {
				SKU: strPtr("5336531800"),
			},
			"2080TI": {
				SKU: strPtr("5218988600"),
			},
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
			"2060": {
				SKU: strPtr("5394903200"),
			},
			"2070": {
				SKU: strPtr("5394901900"),
			},
			"2080": {
				SKU: strPtr("5336531100"),
			},
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
			"2060": {
				SKU: strPtr("5394902900"),
			},
			"2070": {
				SKU: strPtr("5394901600"),
			},
			"2080": {
				SKU: strPtr("5335703700"),
			},
			"2080TI": {
				SKU: strPtr("5218984600"),
			},
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
			"2070": {
				SKU: strPtr("5379432400"),
			},
			"2080": {
				SKU: strPtr("5334463900"),
			},
			"2080TI": {
				SKU: strPtr("5218984100"),
			},
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
	},
	"GBR": {
		Models: map[string]Model{
			"2060": {
				SKU: strPtr("5394903300"),
			},
			"2070": {
				SKU: strPtr("5394902000"),
			},
			"2080": {
				SKU: strPtr("5336531200"),
			},
			"2080TI": {
				SKU: strPtr("5218985600"),
			},
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
			"2060": {
				SKU: strPtr("5394903300"),
			},
			"2070": {
				SKU: strPtr("5394902000"),
			},
			"2080": {
				SKU: strPtr("5336531200"),
			},
			"2080TI": {
				SKU: strPtr("5218985600"),
			},
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
			"2060": {
				SKU: strPtr("5394903400"),
			},
			"2070": {
				SKU: strPtr("5394902100"),
			},
			"2080": {
				SKU: strPtr("5336532000"),
			},
			"2080TI": {
				SKU: strPtr("5218613900"),
			},
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
			"2060": {
				SKU: strPtr("5394903900"),
			},
			"2070": {
				SKU: strPtr("5394902500"),
			},
			"2080": {
				SKU: strPtr("5336531300"),
			},
			"2080TI": {
				SKU: strPtr("5218986100"),
			},
			"3080": {
				SKU: strPtr("5438798100"),
			},
			"3090": {
				SKU: strPtr("5438761600"),
			},
		},
		Locale:       "sv_se",
		NvidiaLocale: "sv-se",
		Currency:     "SEK",
	},
	"LUX": {
		Models: map[string]Model{
			"2060": {
				SKU: strPtr("5394902700"),
			},
			"2070": {
				SKU: strPtr("5336534300"),
			},
			"2080": {
				SKU: strPtr("5336531500"),
			},
			"2080TI": {
				SKU: strPtr("5218987100"),
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
	"POL": {
		Models: map[string]Model{
			"2060": {
				SKU: strPtr("5394903700"),
			},
			"2070": {
				SKU: strPtr("5394902300"),
			},
			"2080": {
				SKU: strPtr("5336531600"),
			},
			"2080TI": {
				SKU: strPtr("5218987600"),
			},
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
			"3080": {
				SKU: strPtr("5438794300"),
			},
		},
		Locale:       "en_gb",
		NvidiaLocale: "en-gb",
		Currency:     "EUR",
	},
	"ESP": {
		Models: map[string]Model{
			"2060": {
				SKU: strPtr("5394903000"),
			},
			"2070": {
				SKU: strPtr("5394901700"),
			},
			"2080": {
				SKU: strPtr("5336531400"),
			},
			"2080TI": {
				SKU: strPtr("5218986600"),
			},
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
			"2060": {
				SKU: strPtr("5394903600"),
			},
			"2070": {
				SKU: strPtr("5394902600"),
			},
			"2080": {
				SKU: strPtr("5336531700"),
			},
			"2080TI": {
				SKU: strPtr("5218988100"),
			},
			"3080": {
				SKU: strPtr("5438797200"),
			},
			"3090": {
				SKU: strPtr("5438797100"),
			},
		},
		Locale:       "no_no",
		NvidiaLocale: "no-no",
		Currency:     "NOK",
	},
	"NLD": {
		Models: map[string]Model{
			"2060": {
				SKU: strPtr("5394903500"),
			},
			"2070": {
				SKU: strPtr("5394902200"),
			},
			"2080": {
				SKU: strPtr("5336532100"),
			},
			"2080TI": {
				SKU: strPtr("5218614400"),
			},
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

func getToast() (*ToastConfig, error) {
	var err error

	switch runtime.GOOS {
	case "linux":
		return &ToastConfig{"linux"}, nil
	case "windows":
		return &ToastConfig{"windows"}, nil
	case "darwin":
		return &ToastConfig{"darwin"}, nil
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return nil, err
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

//getShields Generates ShieldsConfig for application from environmental variables.
func getShields() (*ShieldsConfig, error) {
	c := ShieldsConfig{}

	p, pOk := os.LookupEnv("PORT")
	if pOk == false {
		return nil, &ConfigError{"Shields", "PORT"}
	}
	c.Port = p

	return &c, nil
}

//Get Generates Configuration for application from environmental variables.
func Get(region string, model string, delay int64, sms bool, discord bool, twitter bool, telegram bool, toast bool, shields bool) (*Config, error) {
	if regionConfig, ok := RegionalConfigs[region]; ok {
		models := getSupportedModels(RegionalConfigs[region])
		isSupportedModel := contains(models, model)
		if isSupportedModel == false {
			log.Println(fmt.Sprintf("Please choose one of the following supported models: %v by using -model=XXX", models))
			return nil, &ModelError{"unsupported model error"}
		}
		configuration := Config{}
		configuration.SKU = regionConfig.Models[model].SKU
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

		if toast == true {
			cfg, err := getToast()
			if err != nil {
				return nil, err
			}
			configuration.ToastConfig = cfg
		}

		if shields == true {
			cfg, err := getShields()
			if err != nil {
				return nil, err
			}
			configuration.ShieldsConfig = cfg
		}

		return &configuration, nil
	}

	log.Println(fmt.Sprintf("Please choose one of the following supported regions: %v by using -region=XXX", getSupportedRegions()))
	return nil, &RegionError{region}
}

// contains Determins if a string exists in a slice of strings.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// getSupportedRegions Gets a list of all supported region names
func getSupportedRegions() []string {
	keys := []string{}

	for k := range RegionalConfigs {
		keys = append(keys, k)
	}

	return keys
}

// getSupportedModels Gets a list of all supported model names in a particular region
func getSupportedModels(config RegionalConfig) []string {
	keys := []string{}

	for k := range config.Models {
		keys = append(keys, k)
	}

	return keys
}

// strPtr generates a pointer version of a string
func strPtr(in string) *string {
	i := in
	return &i
}
