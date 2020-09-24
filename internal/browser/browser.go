package browser

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/ianmarmour/nvidia-clerk/internal/config"
)

//OpenProductPage Automatically adds the item to the current cart.
func OpenProductPage(context context.Context, model string, locale string, test bool) error {
	url := fmt.Sprintf("https://www.nvidia.com/%s/geforce/graphics-cards/30-series/rtx-%s/", locale, model)

	err := chromedp.Run(context, chromedp.Navigate(url))
	if err != nil {
		return err
	}

	return nil
}

// Start Starts the ChromeRD browser session and returns it's context.
func Start(config config.Config) (context.Context, error) {
	var allocCtx, _ = chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("enable-automation", false), chromedp.Flag("headless", false))...)
	ctx, _ := chromedp.NewContext(allocCtx)

	return ctx, nil
}
