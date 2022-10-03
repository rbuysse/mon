package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func checkForPattern(body string, pattern string) (check bool) {
	if strings.Contains(body, pattern) {
		check = true
	} else {
		check = false
	}
	return
}

func getPage(url string) (page string, err error) {

	resp, err := http.Get(url)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	page = string(body)

	return page, err
}

func main() {
	for {
		urlToCheck := "https://example.com"
		pattern := "<title>Exampl Domain"

		resp, status := getPage(urlToCheck)

		if status != nil {
			fmt.Printf("ERROR: %s - %s\n", urlToCheck, status)
		} else {
			patternExists := checkForPattern(resp, pattern)
			if patternExists == true {
			} else {
				fmt.Printf("ERROR: %s - Pattern not found: %s\n", urlToCheck, pattern)
			}

		}

		time.Sleep(5 * time.Second)
	}
}
