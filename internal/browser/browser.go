package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/chromedp/chromedp"
)

type Session struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

var session Session

const locale = "en_us"
const nvidiaAPIKey = "9485fa7b159e42edb08a83bde0d83dia"

//exists Determines if a file exists.
func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//getWindowsChromeLocation Determines different Google Chrome install locations.
func getWindowsChromeLocation() string {
	if exists("C:/Program Files (x86)/Google/Chrome/Application/chrome.exe") == true {
		return "C:/Program Files (x86)/Google/Chrome/Application/chrome.exe"
	} else if exists("C:/Program Files/Google/Chrome/Application/chrome.exe") == true {
		return "C:/Program Files/Google/Chrome/Application/chrome.exe"
	}

	userDir, userDirOk := os.LookupEnv("userprofile")
	if userDirOk == false {
		log.Fatal("Unable to determine Google Chrome install location. userprofile env var not set.")
	}

	if exists(userDir + "/AppData/Local/Google/Chrome/Application/chrome.exe") == true {
		return userDir + "/AppData/Local/Google/Chrome/Application/chrome.exe"
	}

	log.Fatal("Unable to determine Google Chrome install location.")
	return "C:/Program Files (x86)/Google/Chrome/Application/chrome.exe"
}

//updateSession Updates the session variable.
func updateSession(sessionResponse string) {
	jsonErr := json.Unmarshal([]byte(sessionResponse), &session)
	if jsonErr != nil {
		log.Fatal("Unable to unmarshal session token.")
	}
}

//getDebugURL Returns the debug information from Chrome running in developer debug mode.
func getDebugURL() string {
	resp, err := http.Get("http://localhost:9222/json/version")
	if err != nil {
		log.Fatal(err)
	}

	var result map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}
	return result["webSocketDebuggerUrl"].(string)
}

// constructSessionURL Builds the session URL
func constructSessionURL() string {
	baseURL := "https://store.nvidia.com/store/nvidia/SessionToken?format=json"
	localeParam := fmt.Sprintf("&locale=%s", locale)
	apiKeyParam := fmt.Sprintf("&apiKey=%s", nvidiaAPIKey)

	return baseURL + localeParam + apiKeyParam
}

//Checkout Opens customer checkout
func Checkout(context context.Context) {
	checkoutURL := fmt.Sprintf("https://api.digitalriver.com/v1/shoppers/me/carts/active/web-checkout?token=%s", session.AccessToken)

	err := chromedp.Run(context,
		chromedp.Navigate(checkoutURL),
	)
	if err != nil {
		log.Fatal(err)
	}
}

//AddToCart Automatically adds the item to the current cart.
func AddToCart(context context.Context, sku string) {
	baseURL := "https://api.digitalriver.com/v1/shoppers/me/carts/active/line-items?format=json&method=post"
	productIDParam := fmt.Sprintf("&productId=%s", sku)
	tokenParam := fmt.Sprintf("&token=%s", session.AccessToken)
	quantityParam := "&quantity=1"
	cartURL := baseURL + productIDParam + tokenParam + quantityParam

	err := chromedp.Run(context,
		chromedp.Navigate(cartURL),
	)
	if err != nil {
		log.Fatal(err)
	}
}

//StartChromeDebugMode Starts a chrome instance in debug-mode
func StartChromeDebugMode() bool {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("google-chrome", "--remote-debugging-port=9222", "--user-data-dir=remote-profile")
	case "windows":
		cmd = exec.Command(getWindowsChromeLocation(), "--remote-debugging-port=9222", "--user-data-dir=remote-profile")
	case "darwin":
		cmd = exec.Command("open", "-a", "/Applications/Google \\Chrome.app", "--args", "--remote-debugging-port=9222", "--user-data-dir=remote-profile")
	default:
		log.Fatal("unsupported platform")
	}

	return cmd.Start() == nil
}

// StartSession Starts the ChromeRD browser session and returns it's context.
func StartSession() context.Context {
	// create allocator context for use with creating a browser context later
	allocatorContext, _ := chromedp.NewRemoteAllocator(context.Background(), getDebugURL())

	// create context
	ctxt, _ := chromedp.NewContext(allocatorContext)

	var sessionResponse string

	err := chromedp.Run(ctxt,
		chromedp.Navigate(constructSessionURL()),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Text(`body > pre`, &sessionResponse),
	)
	if err != nil {
		log.Fatal(err)
	}

	updateSession(sessionResponse)

	return ctxt
}
