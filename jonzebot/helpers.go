package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getTitle(s string) (string, error) {
	// Make HTTP GET request
	response, err := http.Get(s)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	doc, err := html.Parse(response.Body)
	if err != nil {
		panic("Fail to parse html")
	}
	t, e := traverse(doc)
	t = strings.Replace(t, " on Vimeo", "", -1)
	t = strings.Replace(t, " - YouTube", "", -1)
	return t, e
}

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, error) {
	if isTitleElement(n) {
		return n.FirstChild.Data, nil
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, err := traverse(c)
		if err == nil {
			return result, nil
		}
	}

	return "", errors.New("no title found")
}

func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func checkCombo(m, s string) (out bool) {
	for _, v := range playbill {
		if v.Movie == m && v.SkateVid == s {
			out = true
		}
	}
	return
}

func archiveJSON(fn string, ty interface{}) {
	f, err := os.Create(fn)
	if err != nil {
		return
	}

	defer f.Close()

	arch, err := json.Marshal(ty)
	if err != nil {
		return
	}
	c, err := f.Write(arch)
	if err != nil {
		return
	}
	fmt.Println("bytes: ", c)
}

func unarchiveJSON(fn string, ty interface{}) {
	if fileExists(fn) {
		dat, err := ioutil.ReadFile(fn)
		if err != nil {
			return
		}
		json.Unmarshal(dat, ty)
	}
}

func containsVal(s []string, e string) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}
