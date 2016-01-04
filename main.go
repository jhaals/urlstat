package main

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	urls     = kingpin.Arg("urls", "URLs to monitor").Required().Strings()
	diff     = kingpin.Flag("diff", "Display content diff").Short('d').Bool()
	content  = kingpin.Flag("content", "Check for changes in content").Short('c').Bool()
	interval = kingpin.Flag("interval", "Interval between checks.").Default("1s").OverrideDefaultFromEnvar("INTERVAL").Short('i').Duration()
	timeout  = kingpin.Flag("timeout", "HTTP GET timeout.").Default("1s").OverrideDefaultFromEnvar("TIMEOUT").Short('t').Duration()

	red   = color.New(color.FgRed).SprintFunc()
	green = color.New(color.FgGreen).SprintFunc()
)

type result struct {
	Body  string
	Error string
	Code  int
}

func checkURL(url string) result {
	client := http.Client{Timeout: *timeout}
	resp, err := client.Get(url)
	if err != nil {
		return result{"", err.Error(), 0}
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return result{string(body), "", resp.StatusCode}
}

func diffStrings(string1, string2 string) string {
	d := diffmatchpatch.New()

	d1 := d.DiffMain(string1, string2, true)
	var buff bytes.Buffer
	for _, diff := range d1 {
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			buff.WriteString(green(diff.Text))
		case diffmatchpatch.DiffDelete:
			buff.WriteString(red(diff.Text))
		case diffmatchpatch.DiffEqual:
			buff.WriteString(diff.Text)
		}
	}
	return buff.String()
}

func monitor(url string) {
	initalResult := checkURL(url)

	if initalResult.Error != "" {
		fmt.Println(fmt.Sprintf("Checking %s with %s interval. Inital Error %s", url, *interval, initalResult.Error))
	} else {
		fmt.Println(fmt.Sprintf("Checking %s with %s interval. Inital Status %v", url, *interval, initalResult.Code))
	}

	for {
		time.Sleep(*interval)
		result := checkURL(url)

		if result.Error != initalResult.Error && result.Error != "" {
			log.Println(red(result.Error))
		}
		if result.Code != initalResult.Code {
			log.Println(fmt.Sprintf("%s status code changed from %v to %v", url, initalResult.Code, result.Code))
		}
		if result.Body != initalResult.Body && *content && result.Error == "" && initalResult.Error == "" {
			log.Println(url + " content changed")
			if *diff {
				fmt.Println(diffStrings(initalResult.Body, result.Body))
			}
		}
		initalResult = result
	}
}

func main() {
	kingpin.Version("0.0.1")
	kingpin.CommandLine.Help = "Monitor URLs for state changes"
	kingpin.Parse()
	for _, url := range *urls {
		go monitor(url)
	}
	select {}
}
