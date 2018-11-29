package csdn

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"fmt"
)

func Doc(url string) *goquery.Document {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		panic(fmt.Errorf("status code error:%s %d %s", url, res.StatusCode, res.Status))
	}

	// Load the HTML document
	doc, _ := goquery.NewDocumentFromReader(res.Body)

	return doc
}
