package main

import (
	"log"
	"net/http"

	"golang.org/x/net/html"
)

func crawl(link string, client *http.Client, f func(n *html.Node)) {

	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return
	}
	req.Header.Set("Accept", "text/html")

	res, err := client.Do(req)
	if err != nil {
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

	isVisted := make(map[string]struct{})
	var client http.Client
	ch := make(chan string, maxWorker)
	defer close(ch)

	var f func(n *html.Node)
	f = func(n *html.Node) {
		for c := n; c != nil; c = c.NextSibling {
			if c.FirstChild != nil {
				f(c.FirstChild)
			}
			switch c.Data {
			case "a":
				for _, a := range c.Attr {
					if a.Key == "href" {
						ch <- a.Val
						break
					}
				}
			}
		}
	}

	for _, url := range seedUrls {
		go crawl(url, &client, f)
		isVisted[url] = struct{}{}
	}
	for {
		select {
		case url := <-ch:
			if _, ok := isVisted[url]; !ok {
				go crawl(url, &client, f)
				log.Println(url)
				isVisted[url] = struct{}{}
			}
		}
	}

}
