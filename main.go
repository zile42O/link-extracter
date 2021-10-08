package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

func main() {
	color.Blue("Link extracter - 1.0")
	color.Red("Created with ♥ by Zile42O\n\n")
	color.Green("× Please enter the webiste wich you want to grab links..")
	//---------------------------------------------------------
	err := os.Remove("links.txt")
	if err != nil {
	}
	err = os.Remove("output.txt")
	if err != nil {
	}

	buf := bufio.NewReader(os.Stdin)
	color.Blue("> ")
	sentence, err := buf.ReadBytes('\n')
	if err != nil {
		color.Red("Err:", err)
	}
	start := time.Now()
	result := strings.TrimSpace(string(sentence))
	color.Green("× Getting results from: %s", result)
	color.Blue("Please wait...")
	r, err := http.Get(result)
	if err != nil {
		color.Red("Err:", err)
	}
	defer r.Body.Close()
	f, err := os.OpenFile("links.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		color.Red("Err:", err)
	}
	defer f.Close()

	doc, err := goquery.NewDocument(result)
	if err != nil {
		color.Red("Err:", err)
	}
	doc.Find("body a").Each(func(index int, item *goquery.Selection) {
		linkTag := item
		link, _ := linkTag.Attr("href")
		linkText := linkTag.Text()
		if _, err := f.WriteString(link + "\n"); err != nil {
			color.Red("Err:", err)
		}
		fmt.Printf("Output: #%d: '%s' - '%s'\n", index, linkText, link)
	})
	//---------------------------------------------------------
	color.Green("× Please enter the filter word which you want to search in links..")
	buf = bufio.NewReader(os.Stdin)
	color.Blue("> ")
	sentence, err = buf.ReadBytes('\n')
	if err != nil {
		color.Red("Err:", err)
	}
	color.Blue("Please wait...")
	LinebyLineScan(strings.TrimSpace(string(sentence)))
	t := time.Now()
	elapsed := t.Sub(start)
	color.Green("× Program took: %s", elapsed)
	fmt.Println("\n----------------------------------------------------")
	time.Sleep(5000 * time.Millisecond)
}

func LinebyLineScan(filter string) {
	file, err := os.Open("./links.txt")
	if err != nil {
		color.Red("Err:", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	foundnum := 0
	f, err := os.OpenFile("output.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		color.Red("Err:", err)
	}
	defer f.Close()
	for scanner.Scan() {
		query := `^.*` + filter + `.*$`
		match, _ := regexp.MatchString(query, scanner.Text())
		if match {
			fmt.Println(scanner.Text())
			result := scanner.Text() + "\n"
			if _, err := f.WriteString(result); err != nil {
				color.Red("Err:", err)
			}
			foundnum++
		}
	}
	fmt.Println("\n\x1b[35m", foundnum, "links found with filter:", filter, "\x1b[0m")
	if err := scanner.Err(); err != nil {
		color.Red("Err:", err)
	}
}
