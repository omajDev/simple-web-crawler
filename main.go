package main

import (
	"log"
	"net/http"

	"golang.org/x/net/html"
)

func crawl(link chan string, client http.Client, f func(n *html.Node)) {
	req, err := http.NewRequest(http.MethodGet, <-link, nil)
	if err != nil {
		log.Println("ERROR: Failed to crawl:", err)
		return
	}
	req.Header.Set("Accept", "text/html")
	res, err := client.Do(req)
	if err != nil {
		log.Println("ERROR: Failed to crawl:", link)
		return
	}
	defer res.Body.Close()
	n, err := html.Parse(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	f(n)
}

func main() {
	const maxWorker = 5
	seedUrls := []string{
		"https://en.wikipedia.org/wiki",
		"https://blog.hubspot.com/",
	}
	var client http.Client
	ch := make(chan string)
	defer close(ch)
	var f func(n *html.Node)
	f = func(n *html.Node) {

		for c := n; c != nil; c = c.NextSibling {
			if c.FirstChild != nil {
				f(c.FirstChild)
			}

			if c.Data == "a" {
				for _, a := range c.Attr {
					if a.Key == "href" {
						log.Println(a.Val)
					}
				}
			}
		}
	}
	for i := 0; i < maxWorker; i++ {
		go crawl(ch, client, f)
	}
	for _, url := range seedUrls {
		ch <- url
	}
	select {}
}
