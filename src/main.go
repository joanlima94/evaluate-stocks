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

func fetchEarnings(urlCode string, proventosLinks map[string]string) {
	node, err := fetchUrl(urlCode)
	if err != nil {
		fmt.Printf("Error fetching proventos data for URL %s: %v\n", urlCode, err)
		return
	}

	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" && strings.Contains(attr.Val, "proventos.php?papel=") {
					codeProventos := strings.Split(attr.Val, "=")[1]
					fullURL := "https://www.fundamentus.com.br/" + attr.Val
					proventosLinks[codeProventos] = fullURL
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
}

func main() {

	node, err := fetchUrl(constants.UrlFundamentus)
	if err != nil {
		log.Fatal(err)
	}

	stocks := make(map[string]string)
	extractStockData(node, stocks)

	proventosLinks := make(map[string]string)

	for code, link := range stocks {
		fmt.Printf("Code: %s, Link: %s\n", code, link)
		fetchEarnings(link, proventosLinks)
	}

	fmt.Println("Proventos Links:")
	for _, link := range proventosLinks {
		fmt.Println(link)
	}

}
