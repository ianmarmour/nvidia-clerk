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

func uniqueParam() string {
	sec := time.Now().Unix()
	return fmt.Sprintf("&=%v", sec)
}

//updateSession Updates the session variable.
func updateSession(sessionResponse string) {
	jsonErr := json.Unmarshal([]byte(sessionResponse), &session)
	if jsonErr != nil {
		log.Fatal("Unable to unmarshal session token.")
	}
}

// constructSessionURL Builds the session URL
func constructSessionURL(locale string) string {
	baseURL := "https://store.nvidia.com/store/nvidia/SessionToken?format=json"
	localeParam := fmt.Sprintf("&locale=%s", locale)
	apiKeyParam := fmt.Sprintf("&apiKey=%s", nvidiaAPIKey)

	return baseURL + localeParam + apiKeyParam + uniqueParam()
}

// GetInventoryStatus Retrieves sku inventory information from digitalriver
func GetInventoryStatus(ctx context.Context, sku string, locale string) (*InventoryStatus, error) {
	baseURL := fmt.Sprintf("https://api.digitalriver.com/v1/shoppers/me/products/%s/inventory-status?", sku)
	apiKeyParam := fmt.Sprintf("&apiKey=%s", nvidiaAPIKey)
	localeParam := fmt.Sprintf("&locale=%s", locale)
	stockURL := baseURL + apiKeyParam + localeParam + uniqueParam()

	var stockRequestID network.RequestID

	// Have to establish a network listener here to get raw XML response.
	chromedp.ListenTarget(
		ctx,
		func(event interface{}) {
			switch responseReceivedEvent := event.(type) {
			case *network.EventResponseReceived:
				response := responseReceivedEvent.Response
				if response.URL == stockURL {
					stockRequestID = responseReceivedEvent.RequestID
				}
			}
		},
	)

	var stockResponseBody []byte

	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(stockURL),
		chromedp.Sleep(time.Millisecond*500),
		chromedp.ActionFunc(func(cxt context.Context) error {
			body, err := network.GetResponseBody(stockRequestID).Do(cxt)
			stockResponseBody = body
			return err
		}),
	)

	if err != nil {
		log.Println("Error retrieving inventory status")
		return nil, err
	}

	inventoryStatus := InventoryStatus{}

	xmlErr := xml.Unmarshal(stockResponseBody, &inventoryStatus)
	if xmlErr != nil {
		log.Println("Erorr unmarshalling inventory XML")
		return nil, xmlErr
	}

	return &inventoryStatus, nil
}

//Checkout Opens customer checkout
func Checkout(context context.Context, locale string) error {
	checkoutURL := fmt.Sprintf("https://api.digitalriver.com/v1/shoppers/me/carts/active/web-checkout?token=%s&locale=%s", session.AccessToken, locale) + uniqueParam()

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
	baseURL := "https://api.digitalriver.com/v1/shoppers/me/carts/active/line-items?format=json&method=post"
	productIDParam := fmt.Sprintf("&productId=%s", sku)
	tokenParam := fmt.Sprintf("&token=%s", session.AccessToken)
	quantityParam := "&quantity=1"
	localeParam := fmt.Sprintf("&locale=%s", locale)
	cartURL := baseURL + productIDParam + tokenParam + quantityParam + localeParam + uniqueParam()

	err := chromedp.Run(context,
		chromedp.Navigate(cartURL),
	)
	if err != nil {
		return err
	}

	return nil
}

// StartSession Starts the ChromeRD browser session and returns it's context.
func StartSession(config config.Config) context.Context {
	// create allocator context for use with creating a browser context later

	allocatorContext, _ := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)

	// create context
	ctxt, _ := chromedp.NewContext(allocatorContext)

	var sessionResponse string

	err := chromedp.Run(ctxt,
		chromedp.Navigate(constructSessionURL(config.Locale)),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Text(`body > pre`, &sessionResponse),
	)
	if err != nil {
		log.Fatal(err)
	}

	updateSession(sessionResponse)

	return ctxt
}
