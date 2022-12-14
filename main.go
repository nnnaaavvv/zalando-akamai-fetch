package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var url = "ikea.no"

func FetchParse(url string) string {
	resp, err := http.Get("https://" + url)
	if err != nil {
		panic(err)
	}
	/*
		html, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
	*/

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// checks html for tag <script type="text/javascript" src=
	var uri []string
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		if typee, _ := s.Attr("type"); typee == "text/javascript" {
			textjavascript, _ := s.Attr("src")
			uri = append(uri, textjavascript)
		}
	})
	if len(uri) == 0 {
		fmt.Println("link not found")
		os.Exit(1)
	}
	for _, uri := range uri {
		fullURL := uri
		//WfullURL := "https://" + url + uri
		if strings.Contains(uri, "https://") == false {
			response, err := http.Get("https://" + url + fullURL)
			if err != nil {
				panic(err)
			}

			html, err := io.ReadAll(response.Body)
			if err != nil {
				panic(err)
			}

			if strings.Contains(string(html), "bmak") == true {
				fmt.Println("https://" + url + fullURL)
				break
			} else {
				continue
			}
		}
		response, err := http.Get(fullURL)
		if err != nil {
			panic(err)
		}

		html, err := io.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		if strings.Contains(string(html), "bmak") == true {
			fmt.Println(fullURL)
			break
		} else {
			continue
		}
	}
	var CompleteURL string
	if strings.Contains(uri[len(uri)-1], "https://") == false {
		CompleteURL = "https://" + url + uri[len(uri)-1]
	} else if strings.Contains(uri[len(uri)-1], "https://") == true {
		CompleteURL = uri[len(uri)-1]
	}

	script, err := http.Get(CompleteURL)
	if err != nil {
		panic(err)
	}

	html, err := io.ReadAll(script.Body)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("obfuscated.js")
	defer f.Close()

	f.Write(html)
	return CompleteURL
}

func main() {
	start := time.Now()
	FetchParse(url)
	elapsedTime := time.Since(start)
	log.Printf("%s", elapsedTime)
}
