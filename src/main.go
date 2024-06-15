package main

import (
	"evaluate-stocks/constants"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func fetchUrl(url string) (*html.Node, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	node, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	return node, err

}

func extractStockData(node *html.Node, stocks map[string]string) {
	if node.Type == html.ElementNode && node.Data == "a" {

		for _, attr := range node.Attr {
			if attr.Key == "href" && strings.Contains(attr.Val, "detalhes.php?papel=") {
				code := strings.Split(attr.Val, "=")[1]
				stocks[code] = "https://www.fundamentus.com.br/" + attr.Val
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		extractStockData(c, stocks)
	}

}

func main() {

	node, err := fetchUrl(constants.UrlFundamentus)
	if err != nil {
		log.Fatal(err)
	}

	stocks := make(map[string]string)
	extractStockData(node, stocks)

	for code, url := range stocks {
		fmt.Printf("Code: %s, Link: %s\n", code, url)
	}

}
