package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func f(n *html.Node) {

	for c := n; c != nil; c = c.NextSibling {
		if c.FirstChild != nil {
			f(c.FirstChild)
		}

		if c.Data == "a" {
			for _, a := range c.Attr {
				if a.Key == "href" {
					fmt.Println(a.Val)
				}
			}
		}
	}
}
func crawl(link string) {

	resp, err := http.Get(link)
	if err != nil {
		fmt.Println("ERROR: Failed to crawl:", link)
		return
	}
	defer resp.Body.Close()
	n, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	f(n)
}

func main() {
	seedUrls := []string{
		"https://en.wikipedia.org/wiki",
	}
	ch := make(chan string)
	defer close(ch)
	for _, url := range seedUrls {
		crawl(url)
	}
}
