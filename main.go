package main

import (
	"context"
	"errors"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
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
	Url      string
	Error    error
	Duration float64
}

func checkForPattern(body string, pattern string) (check bool) {
	if strings.Contains(body, pattern) {
		check = true
	} else {
		check = false
	}
	return
}

func dbConnect() (influxdb2.Client, error) {

	token := os.Getenv("INFLUXDB_TOKEN")
	if token == "" {
		return nil, errors.New("INFLUXDB_TOKEN must be set")
	}

	url := os.Getenv("INFLUXDB_URL")
	if url == "" {
		return nil, errors.New("INFLUXDB_URL must be set")
	}

	client := influxdb2.NewClient(url, token)
	_, err := client.Health(context.Background())

	return client, err
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
	prefetch := time.Now()
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
			m := fmt.Sprintf("ERROR: %s - Pattern not found: %s", url, pattern)
			err := errors.New(m)
			res.Error = err
		}
	}

	res.Duration = time.Since(prefetch).Seconds()
	ch <- res
}

func main() {
	configFile := "config.toml"
	for {
		c := getConfig(configFile)

		const bucket = "mon"
		const org = "mon"

		client, err := dbConnect()

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}

		writeAPI := client.WriteAPIBlocking("mon", "mon")

		loopstart := time.Now()
		ch := make(chan Result)

		for i := 0; i < len(c.Sites); i++ {
			urlToCheck := c.Sites[i].Url
			pattern := c.Sites[i].Pattern
			go checkPage(urlToCheck, pattern, ch)

		}

		for i := 0; i < len(c.Sites); i++ {
			var dbErr string
			ret := <-ch
			if ret.Error != nil {
				dbErr = fmt.Sprintf("%v", ret.Error)
			}
			p := influxdb2.NewPointWithMeasurement("stat").
				AddTag("host", ret.Url).
				AddField("resp", ret.Duration).
				AddField("error", dbErr).
				SetTime(time.Now())
			writeAPI.WritePoint(context.Background(), p)
			fmt.Printf("%s - %v - %.2fs\n", ret.Url, ret.Error, ret.Duration)
		}
		fmt.Printf("%.2fs elapsed\n\n", time.Since(loopstart).Seconds())
		time.Sleep(5 * time.Second)
	}
}
