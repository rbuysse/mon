package main

import (
	"errors"
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

type Result struct {
	Url   string
	Error error
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

func checkPage(url string, pattern string, ch chan<- Result) {
	var res Result
	res.Url = url
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if err != nil {
		res.Error = err
	} else {
		patternExists := checkForPattern(string(body), pattern)
		if patternExists == true {
		} else {
			m := fmt.Sprintf("ERROR: %s - Pattern not found: %s\n", url, pattern)
			err := errors.New(m)
			res.Error = err
		}
	}

	ch <- res
}

func main() {
	configFile := "config.toml"
	for {
		c := getConfig(configFile)
		ch := make(chan Result)

		for i := 0; i < len(c.Sites); i++ {
			urlToCheck := c.Sites[i].Url
			pattern := c.Sites[i].Pattern
			go checkPage(urlToCheck, pattern, ch)

		}

		for i := 0; i < len(c.Sites); i++ {
			ret := <-ch
			fmt.Println(ret.Url)
		}
		time.Sleep(5 * time.Second)
	}
}
