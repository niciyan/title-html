package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: title-html <URL>")
		return
	}

	for i, url := range os.Args[1:] {
		if i >= 1 {
			fmt.Println("")
		}
		err := fetchTitle(url)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func fetchTitle(url string) error {

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error while fetching %v. %v\n", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 400 {
		return errors.New(fmt.Sprintf("Request failed. Response Status Code: %v", res.StatusCode))
	}

	node, err := html.Parse(res.Body)
	if err != nil {
		return fmt.Errorf("error while body. %v", err)
	}
	res.Body.Close()

	n, err := getTitleNode(node)
	if err != nil {
		return err
	}
	fmt.Printf("[%v](%v)\n", n.FirstChild.Data, url)
	return nil
}

func getTitleNode(n *html.Node) (*html.Node, error) {
	var ret *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			ret = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	if ret != nil {
		return ret, nil
	}
	return nil, errors.New("Missing <title> in the html tree")
}
