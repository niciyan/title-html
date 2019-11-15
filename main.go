package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: title-html <URL>")
		return
	}

	url := os.Args[1]

	res, err := http.Get(url)
	if err != nil {
		fmt.Errorf("error while fetching %v. %v\n", url, err)
		return
	}
	defer res.Body.Close()

	err = handleStatus(res)
	if err != nil {
		fmt.Println(err)
		return
	}

	node, err := html.Parse(res.Body)
	if err != nil {
		fmt.Errorf("error while body. %v", err)
		return
	}
	res.Body.Close()

	n, err := getTitleNode(node)
	fmt.Printf("[%v](%v)\n", n.FirstChild.Data, url)
	// title, _ := traverse(node)
	// fmt.Printf("[%v](%v)\n", title, url)
}

func handleStatus(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode > 400 {
		return errors.New(fmt.Sprintf("Request seems to fail. Response Status Code: %v", resp.StatusCode))
	}
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
