package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// ProductsResponse Used for unmarshalling JSON response from api.nvidia.partners
type ProductsResponse struct {
	Products Products `json:"products"`
}

// Products Used for unmarshalling JSON response from api.nvidia.partners
type Products struct {
	Product []Product `json:"product"`
}

type Product struct {
	ID                     int64            `json:"id"`
	Name                   string           `json:"name"`
	DisplayName            string           `json:"displayName"`
	SKU                    string           `json:"sku"`
	DisplayableProduct     string           `json:"displayableProduct"`
	ManufacturerPartNumber string           `json:"manufacturerPartNumber"`
	MaximumQuantity        int64            `json:"maximumQuantity"`
	ThumbnailImage         string           `json:"thumbnailImage"`
	CustomAttributes       CustomAttributes `json:"customAttributes"`
	Pricing                Pricing          `json:"pricing"`
	InventoryStatus        InventoryStatus  `json:"inventoryStatus"`
	RelatedProducts        RelatedProducts  `json:"relatedProducts"`
	ViewStyle              string           `json:"viewStyle"`
}

type CustomAttributes struct {
	Attribute []Attribute `json:"attribute"`
}

type Attribute struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Pricing struct {
	URI                                string   `json:"uri"`
	ListPrice                          Price    `json:"listPrice"`
	ListPriceWithQuantity              Price    `json:"listPriceWithQuantity"`
	SalePriceWithQuantity              Price    `json:"salePriceWithQuantity"`
	FormattedListPrice                 string   `json:"formattedListPrice"`
	FormattedListPriceWithQuantity     string   `json:"formattedListPriceWithQuantity"`
	FormattedSalePriceWithQuantity     string   `json:"formattedSalePriceWithQuantity"`
	TotalDiscountWithQuantity          Price    `json:"totalDiscountWithQuantity"`
	FormattedTotalDiscountWithQuantity string   `json:"formattedTotalDiscountWithQuantity"`
	ListPriceIncludesTax               string   `json:"listPriceIncludesTax"`
	Tax                                Tax      `json:"tax"`
	FeePricing                         FeePrice `json:"feePricing"`
}

type Price struct {
	Currency string `json:"currency"`
	Value    int64  `json:"value"`
}

type Tax struct {
	VatPercentage int64 `json:"vatPercentage"`
}

type FeePrice struct {
	SalePriceWithFeesAndQuantity          Price  `json:"salePriceWithFeesAndQuantity"`
	FormattedSalePriceWithFeesAndQuantity string `json:"formattedSalePriceWithFeesAndQuantity"`
}

type InventoryStatus struct {
	URI                          string `json:"uri"`
	AvailableQuantityIsEstimated string `json:"availableQuantityIsEstimated"`
	ProductIsInStock             string `json:"productIsInStock"`
	ProductIsAllowsBackorders    string `json:"productIsAllowsBackorders"`
	ProductIsTracked             string `json:"productIsTracked"`
	RequestedQuantityAvailable   string `json:"requestedQuantityAvailable"`
	Status                       string `json:"status"`
	StatusIsEstimated            string `json:"statusIsEstimated"`
}

type RelatedProducts []RelatedProduct

type RelatedProduct struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Inventory struct {
	Product InventoryProduct `json:"Product"`
}

type InventoryProduct struct {
	AvailableQuantity int64 `json:"availableQuantity"`
}

type SessionToken struct {
	Value string `json:"session_token"`
}

type AddToCartResponse struct {
	URL string `json:"location"`
}

//GetSessionToken Retrieves the session token for NVIDIA store.
func GetSessionToken(client *http.Client) (*SessionToken, error) {
	url := "https://store.nvidia.com/store/nvidia/SessionToken?format=json" + urlTime()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", getUserAgent())

	resBody, err := getBody(req, client)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	session := SessionToken{}
	jsonErr := json.Unmarshal(resBody, &session)
	if jsonErr != nil {
		log.Println(jsonErr)
		return nil, jsonErr
	}

	return &session, nil
}

//AddToCheckout Adds a product to a checkout cart for NVIDIA store.
func AddToCheckout(sku string, token string, locale string, client *http.Client) (*AddToCartResponse, error) {
	url := "https://api-prod.nvidia.com/direct-sales-shop/DR/add-to-cart"
	reqBody := []byte(fmt.Sprintf(`{"products": [{"productId":%s,"quantity": 1}]}`, sku))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("locale", locale)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("nvidia_shop_id", token)
	req.Header.Add("charset", "utf-8")

	resBody, err := getBody(req, client)
	if err != nil {
		return nil, err
	}

	cart := AddToCartResponse{}

	jsonErr := json.Unmarshal(resBody, &cart)
	if jsonErr != nil {
		return nil, jsonErr
	}

	cart.URL = cart.URL + urlTime()

	return &cart, nil
}

// GetSkuInfo Looks up SKU invormation from NVIDIAs API.
func GetSkuInfo(sku string, locale string, currency string, client *http.Client) (*ProductsResponse, error) {
	url := fmt.Sprintf("https://api-prod.nvidia.com/direct-sales-shop/DR/products/%s/%s/%s", locale, currency, sku)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resBody, err := getBody(req, client)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := ProductsResponse{}
	jsonErr := json.Unmarshal(resBody, &products)
	if jsonErr != nil {
		log.Println(jsonErr)
		return nil, jsonErr
	}

	return &products, nil
}

func getUserAgent() string {
	agent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36"
	return agent
}

// Gets the byte data out of a request body.
func getBody(request *http.Request, client *http.Client) ([]byte, error) {
	r, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if r.StatusCode > 400 {
		return nil, err
	}
	defer r.Body.Close()

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		return nil, readErr
	}

	return body, nil
}

//urlTime Generates a url encoded datetime parameter used for cache invalidation in browsers.
func urlTime() string {
	sec := time.Now().Unix()
	return fmt.Sprintf("&=%v", sec)
}
