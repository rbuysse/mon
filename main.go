package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Sites []*Site `toml:"site"`
}

type Site struct {
	Url     string `toml:"url"`
	Pattern string `toml:"pattern"`
}

func checkForPattern(body string, pattern string) (check bool) {
	if strings.Contains(body, pattern) {
		check = true
	} else {
		check = false
	}
	return
}

func getConfig(f string) (c Config) {
	if _, err := os.Stat(f); err != nil {
		log.Fatal(err)
	}

	_, err := toml.DecodeFile(f, &c)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return c
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
	configFile := "config.toml"
	for {
		c := getConfig(configFile)
		for i := 0; i < len(c.Sites); i++ {
			fmt.Println(c.Sites[i].Url)
			urlToCheck := c.Sites[i].Url
			pattern := c.Sites[i].Pattern

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
		}
		time.Sleep(5 * time.Second)
	}
}
