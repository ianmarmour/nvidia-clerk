package browser

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

type InventoryStatus struct {
	XMLName                      xml.Name `xml:"inventoryStatus"`
	URI                          string   `xml:"uri,attr"`
	Product                      Product  `xml:"product"`
	AvailableQuantityIsEstimated bool     `xml:"availableQuantityIsEstimated"`
	ProductIsInStock             bool     `xml:"productIsInStock"`
	ProductIsAllowsBackorders    bool     `xml:"productIsAllowsBackorders"`
	ProductIsTracked             bool     `xml:"productIsTracked"`
	RequestedQuantityAvailable   bool     `xml:"requestedQuantityAvailable"`
	Status                       string   `xml:"status"`
	StatusIsEstimated            bool     `xml:"statusIsEstimated"`
	CustomStockMessage           string   `xml:"customStockMessage"`
}

type Product struct {
	XMLName             xml.Name `xml:"product"`
	URI                 string   `xml:"uri,attr"`
	ID                  string   `xml:"id"`
	ExternalReferenceID string   `xml:"externalReferenceId"`
	CompanyID           string   `xml:"companyId"`
}

type Session struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

var session Session

const nvidiaAPIKey = "9485fa7b159e42edb08a83bde0d83dia"

//urlTime Generates a url encoded datetime parameter used for cache invalidation in browsers.
func urlTime() string {
	sec := time.Now().Unix()
	return fmt.Sprintf("&=%v", sec)
}

//updateSession Updates the session variable.
func updateSession(sessionResponse string) error {
	err := json.Unmarshal([]byte(sessionResponse), &session)
	if err != nil {
		return err
	}

	return nil
}

// sessionURL Constructs the url for the NVIDIA store to start your session
func sessionURL(locale string) string {
	url := "https://store.nvidia.com/store/nvidia/SessionToken?format=json"
	loc := fmt.Sprintf("&locale=%s", locale)
	api := fmt.Sprintf("&apiKey=%s", nvidiaAPIKey)

	return url + loc + api + urlTime()
}

// GetInventoryStatus Retrieves sku inventory information from digitalriver
func GetInventoryStatus(ctx context.Context, sku string, locale string, delay int64) (*InventoryStatus, error) {
	url := fmt.Sprintf("https://api.digitalriver.com/v1/shoppers/me/products/%s/inventory-status?", sku)
	api := fmt.Sprintf("&apiKey=%s", nvidiaAPIKey)
	loc := fmt.Sprintf("&locale=%s", locale)
	url = url + api + loc + urlTime()

	// Avoid race conditions in ActionFunc
	reqIDs := make(chan network.RequestID)

	// Have to establish a network listener here to get raw XML response.
	chromedp.ListenTarget(
		ctx,
		func(event interface{}) {
			switch ev := event.(type) {
			case *network.EventResponseReceived:
				go func() {
					response := ev.Response
					if response.URL == url {
						reqIDs <- ev.RequestID
					}
				}()
			}
		},
	)

	var resBody []byte

	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(url),
		chromedp.Sleep(time.Millisecond*time.Duration(delay)),
		chromedp.ActionFunc(func(cxt context.Context) error {
			id := <-reqIDs
			body, err := network.GetResponseBody(id).Do(cxt)
			resBody = body
			return err
		}),
	)
	if err != nil {
		log.Println("Error retrieving inventory status")
		return nil, err
	}

	inventoryStatus := InventoryStatus{}

	xmlErr := xml.Unmarshal(resBody, &inventoryStatus)
	if xmlErr != nil {
		log.Println("Erorr unmarshalling inventory XML")
		return nil, xmlErr
	}

	return &inventoryStatus, nil
}

//Checkout Opens customer checkout
func Checkout(context context.Context, locale string) error {
	checkoutURL := fmt.Sprintf("https://api.digitalriver.com/v1/shoppers/me/carts/active/web-checkout?token=%s&locale=%s", session.AccessToken, locale) + urlTime()

	err := chromedp.Run(context,
		chromedp.Navigate(checkoutURL),
	)
	if err != nil {
		return err
	}

	return nil
}

//AddToCart Automatically adds the item to the current cart.
func AddToCart(context context.Context, sku string, locale string) error {
	url := "https://api.digitalriver.com/v1/shoppers/me/carts/active/line-items?format=json&method=post"
	id := fmt.Sprintf("&productId=%s", sku)
	tok := fmt.Sprintf("&token=%s", session.AccessToken)
	qty := "&quantity=1"
	loc := fmt.Sprintf("&locale=%s", locale)
	url = url + id + tok + qty + loc + urlTime()

	err := chromedp.Run(context, chromedp.Navigate(url))
	if err != nil {
		return err
	}

	return nil
}

// Start Starts the ChromeRD browser session and returns it's context.
func Start(config config.Config) (context.Context, error) {
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
	ctx, _ := chromedp.NewContext(allocCtx)
	var res string

	err := chromedp.Run(ctx,
		chromedp.Navigate(sessionURL(config.Locale)),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Text(`body > pre`, &res),
	)
	if err != nil {
		return nil, err
	}

	err = updateSession(res)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}
