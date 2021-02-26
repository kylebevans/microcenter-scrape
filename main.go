package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func main() {

	client := &http.Client{}

	var pages [4]string
	pages[0] = "https://www.microcenter.com/search/search_results.aspx?Ntt=rtx+3080&Ntk=all&sortby=match&N=0&myStore=true"
	pages[1] = "https://www.microcenter.com/search/search_results.aspx?Ntt=rtx+3070&Ntk=all&sortby=match&N=0&myStore=true"
	pages[2] = "https://www.microcenter.com/search/search_results.aspx?Ntt=rx+6800&Ntk=all&sortby=match&N=0&myStore=true"
	pages[3] = "https://www.microcenter.com/search/search_results.aspx?Ntt=rx+6900&Ntk=all&sortby=match&N=0&myStore=true"
	for _, page := range pages {
		req, err := http.NewRequest("GET", page, nil)
		if err != nil {
			panic(err)
		}

		req.Header.Set("Cookie", "storeSelected=141")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		doc, err := html.Parse(resp.Body)
		if err != nil {
			panic(err)
		}

		var f func(*html.Node, bool, *bool, *[]string)
		f = func(n *html.Node, inDetail bool, inStock *bool, productNames *[]string) {
			if inDetail && n.Type == html.ElementNode && n.Data == "a" {
				for _, v := range n.Attr {
					if v.Key == "data-name" {
						*productNames = append(*productNames, v.Val+"\n")
					}
				}
			}
			// if inDetail && len(n.Attr) > 0 && n.Attr[0].Val == "stock" {
			// 	// if n.FirstChild.FirstChild.FirstChild.Data != "SOLD OUT" {
			// 	// 	fmt.Println(desc)
			// 	// 	return
			// 	// }
			// 	fmt.Println(n.FirstChild.Type)
			// }
			if n.Type == html.ElementNode && len(n.Attr) > 0 && n.Attr[0].Val == "detail_wrapper" {
				inDetail = true
			}
			if n.Type == html.ElementNode && len(n.Attr) > 0 && n.Attr[0].Val == "instore-nostock " {
				*inStock = false
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c, inDetail, inStock, productNames)
			}
		}
		var productNames []string
		inStock := true
		f(doc, false, &inStock, &productNames)
		if inStock {
			fmt.Println(productNames)
		}

	}

}
