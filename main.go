package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func main() {
	var simple = flag.Bool("s", false, "simple mode")
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: title-html <URL> <URL> ...")
		return
	}

	for _, url := range flag.Args() {
		res, err := http.Get(url)
		if err != nil {
			fmt.Printf("error while fetching %v. %v\n", url, err)
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
			fmt.Printf("error while body. %v", err)
			return
		}
		res.Body.Close()

		n, err := getTitleNode(node)
		if *simple {
			fmt.Printf("%v\n%v\n", n.FirstChild.Data, url)
		} else {
			fmt.Printf("[%v](%v)\n", n.FirstChild.Data, url)
		}

	}
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
