package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ProductsResponse struct {
	Products Products `json:products`
}

type Products struct {
	Product []Product `json:"product"`
}

type Product struct {
	ID                     int              `json:"id"`
	Name                   string           `json:"name"`
	DisplayName            string           `json:"displayName"`
	SKU                    string           `json:"sku"`
	DisplayableProduct     string           `json:"displayableProduct"`
	ManufacturerPartNumber string           `json:"manufacturerPartNumber"`
	MaximumQuantity        int              `json:"maximumQuantity"`
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
	Value    int    `json:"value"`
}

type Tax struct {
	VatPercentage int `json:"vatPercentage"`
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
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Inventory struct {
	Product InventoryProduct `json:"Product"`
}

type InventoryProduct struct {
	AvailableQuantity int `json:"availableQuantity"`
}

func GetSkuInfo(sku string, locale string, currency string, client *http.Client) (*ProductsResponse, error) {
	endpoint := fmt.Sprintf("https://api-prod.nvidia.com/direct-sales-shop/DR/products/%s/%s/%s", locale, currency, sku)

	r, err := client.Get(endpoint)
	if err != nil {
		message := fmt.Sprintf("Error attempting to access URL: %s", endpoint)
		log.Println(message)
		return nil, err
	}

	if r.Body != nil {
		defer r.Body.Close()
	}

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Println(readErr)
		return nil, readErr
	}

	products := ProductsResponse{}
	jsonErr := json.Unmarshal(body, &products)
	if jsonErr != nil {
		log.Println(jsonErr)
		return nil, jsonErr
	}

	return &products, nil
}
