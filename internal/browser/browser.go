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

type Products struct {
	XMLName      xml.Name      `xml:"products"`
	URI          string        `xml:"uri,attr"`
	NextPage     *NextPage     `xml:"nextPage,omitempty"`
	PreviousPage *PreviousPage `xml:"previousPage,omitempty"`
	Product      []Product     `xml:"product"`
}

type PreviousPage struct {
	XMLName xml.Name `xml:"previousPage"`
	URI     string   `xml:"uri,attr"`
}

type NextPage struct {
	XMLName xml.Name `xml:"nextPage"`
	URI     string   `xml:"uri,attr"`
}

// Product Digital River product schema used for XML unmarshalling
type Product struct {
	XMLName                xml.Name          `xml:"product"`
	URI                    string            `xml:"uri,attr,omitempty"`
	Categories             *Categories       `xml:"categories,omitempty"`
	FamilyAttributes       *FamilyAttributes `xml:"familyAttributes,omitempty"`
	ID                     string            `xml:"id,omitempty"`
	Name                   string            `xml:"name,omitempty"`
	DisplayName            string            `xml:"displayName,omitempty"`
	ShortDescription       string            `xml:"shortDescription,omitempty"`
	LongDescription        string            `xml:"longDescription,omitempty"`
	ProductType            string            `xml:"productType,omitempty"`
	SKU                    string            `xml:"sku,omitempty"`
	ExternalReferenceID    string            `xml:"externalReferenceId,omitempty"`
	CompanyID              string            `xml:"companyId,omitempty"`
	DisplayableProduct     bool              `xml:"displayableProduct,omitempty"`
	Purchasable            bool              `xml:"purchasable,omitempty"`
	ManufacturerName       string            `xml:"manufacturerName,omitempty"`
	ManufacturerPartNumber string            `xml:"manufacturerPartNumber,omitempty"`
	MinimumQuantity        int               `xml:"minimumQuantity,omitempty"`
	MaximumQuantity        int               `xml:"maximumQuantity,omitempty"`
	ThumbnailImage         string            `xml:"thumbnailImage,omitempty"`
	ProductImage           string            `xml:"productImage,omitempty"`
	Keywords               string            `xml:"keywords,omitempty"`
	BaseProduct            bool              `xml:"baseProduct,omitempty"`
	Pricing                *ProductPricing   `xml:"pricing,omitempty"`
	AddProductToCart       *AddProductToCart `xml:"addProductToCart,omitempty"`
}

// ProductPricing Digital River pricing schema used for XML unmarshalling of Product
type ProductPricing struct {
	XMLName                        xml.Name               `xml:"pricing"`
	URI                            string                 `xml:"uri,attr,omitempty"`
	ListPrice                      *ListPrice             `xml:"listPrice,omitempty"`
	SalePriceWithQuantity          *SalePriceWithQuantity `xml:"salePriceWithQuantity,omitempty"`
	FormattedListPrice             string                 `xml:"formattedListPrice,omitempty"`
	FormattedSalePriceWithQuantity string                 `xml:"formattedSalePriceWithQuantity,omitempty"`
	ListPriceIncludesTax           bool                   `xml:"listPriceIncludesTax,omitempty"`
	MSRPPrice                      string                 `xml:"msrpPrice,omitempty"`
	FormattedMSRPPrice             string                 `xml:"formattedMsrpPrice,omitempty"`
}

// Categories Digital River categories schema used for XML unmarshalling of Product
type Categories struct {
	XMLName xml.Name `xml:"categories"`
	URI     string   `xml:"uri,attr,omitempty"`
}

// FamilyAttributes Digital River familyAttributes schema used for XML unmarshalling of Product
type FamilyAttributes struct {
	XMLName xml.Name `xml:"familyAttributes"`
	URI     string   `xml:"uri,attr,omitempty"`
}

// ListPrice Digital River listPrice schema used for XML unmarshalling of Product
type ListPrice struct {
	XMLName  xml.Name `xml:"listPrice"`
	Currency string   `xml:"currency,attr,omitempty"`
}

// SalePriceWithQuantity Digital River salePriceWithQuantity schema used for XML unmarshalling of Product
type SalePriceWithQuantity struct {
	XMLName  xml.Name `xml:"salePriceWithQuantity"`
	Currency string   `xml:"currency,attr,omitempty"`
}

//AddProductToCart Digital River addProductToCart schema used for XML unmarshalling of Product
type AddProductToCart struct {
	XMLName xml.Name `xml:"addProductToCart"`
	CartURI string   `xml:"cartUri,attr,omitempty"`
	URI     string   `xml:"uri,attr,omitempty"`
}

//InventoryStatus Digital River inventoryStatus schema used for XML unmarshalling
type InventoryStatus struct {
	XMLName                      xml.Name         `xml:"inventoryStatus"`
	URI                          string           `xml:"uri,attr"`
	Product                      InventoryProduct `xml:"product"`
	AvailableQuantityIsEstimated bool             `xml:"availableQuantityIsEstimated"`
	ProductIsInStock             bool             `xml:"productIsInStock"`
	ProductIsAllowsBackorders    bool             `xml:"productIsAllowsBackorders"`
	ProductIsTracked             bool             `xml:"productIsTracked"`
	RequestedQuantityAvailable   bool             `xml:"requestedQuantityAvailable"`
	Status                       string           `xml:"status"`
	StatusIsEstimated            bool             `xml:"statusIsEstimated"`
	CustomStockMessage           string           `xml:"customStockMessage"`
}

//InventoryProduct Digital River product schema used for XML unmarshalling of inventoryStatus
type InventoryProduct struct {
	XMLName             xml.Name `xml:"product"`
	URI                 string   `xml:"uri,attr"`
	ID                  string   `xml:"id"`
	ExternalReferenceID string   `xml:"externalReferenceId"`
	CompanyID           string   `xml:"companyId"`
}

//Session NVIDIA store session response used to authenticate in the browser
type Session struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

var session Session
var allocCtx, _ = chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("enable-automation", false), chromedp.Flag("headless", false))...)

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

// GetProducts Returns a list of all products in the Digital River Inventory
func GetProducts(ctx context.Context, locale string, delay int64) ([]*Product, error) {
	pageNum := 1
	products := []*Product{}

	for {
		current, err := GetProductsPage(ctx, locale, delay, pageNum)
		if err != nil {
			return nil, err
		}

		for i := range current.Product {
			products = append(products, &current.Product[i])
		}

		if current.NextPage == nil {
			break
		}

		pageNum++
	}

	return products, nil
}

// GetProductsPage Retrieves a page of products from Digital River
func GetProductsPage(ctx context.Context, locale string, delay int64, page int) (*Products, error) {
	url := "https://api.digitalriver.com/v1/shoppers/me/products?"
	pag := fmt.Sprintf("&pageNumber=%v", page)
	api := fmt.Sprintf("&apiKey=%s", nvidiaAPIKey)
	loc := fmt.Sprintf("&locale=%s", locale)
	fil := fmt.Sprintf("&fields=%s", "product.id,product.displayName")
	url = url + pag + api + loc + fil + urlTime()

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

	products := Products{}

	xmlErr := xml.Unmarshal(resBody, &products)
	if xmlErr != nil {
		log.Println("Erorr unmarshalling inventory XML")
		return nil, xmlErr
	}

	return &products, nil
}

// GetProduct Retrieves sku product information from digitalriver
func GetProduct(ctx context.Context, sku string, locale string, delay int64) (*Product, error) {
	url := fmt.Sprintf("https://api.digitalriver.com/v1/shoppers/me/products/%s?", sku)
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

	product := Product{}

	xmlErr := xml.Unmarshal(resBody, &product)
	if xmlErr != nil {
		log.Println("Erorr unmarshalling inventory XML")
		return nil, xmlErr
	}

	return &product, nil
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
