package href

import (
	"fmt"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"time"
)

func Get() (string, error) {
	hrefCh := make(chan string, 1)
	go func() {
		geziyor.NewGeziyor(&geziyor.Options{
			StartURLs: []string{"https://urpet96.ru/?page_id=11619"},
			ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
				if href, ok := r.HTMLDoc.Find("a.et_pb_button.et_pb_promo_button").Attr("href"); ok {
					hrefCh <- href
				}
			},
			LogDisabled: true,
		}).Start()
	}()

	select {
	case href := <-hrefCh:
		return href, nil
	case <-time.After(5 * time.Second):
		return "", fmt.Errorf("timeout after 5 seconds")
	}
}
