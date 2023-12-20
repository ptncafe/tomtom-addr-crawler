package main

import (
	"encoding/csv"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"io"
	"log"
	"os"
)

const domain_url = "https://plan.tomtom.com/en/?p=10.82734,106.66315,9.55z&q=10.76397248,106.6881186"

func main() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	if _, err = page.Goto(domain_url); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	f, err := os.Open("data/data.csv")
	if err != nil {
		log.Fatal(err)
	}
	fileResult, err := os.Create("data/result.csv")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()
	defer fileResult.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	writer := csv.NewWriter(fileResult)
	defer writer.Flush()
	index := -1
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		index++
		if index <= 0 {
			continue
		}
		// do something with read line
		log.Printf("%+v\n", rec)
		if err = page.Locator("input#search").Click(); err != nil {
			log.Fatalf("(\"input#search\").Click(): %v", err)
		}
		if err = page.Locator("input#search").Clear(); err != nil {
			log.Fatalf("(\"input#search\").Click(): %v", err)
		}
		latLong := fmt.Sprintf("%s,%s", rec[1], rec[2])
		err = page.Locator("input#search").Fill(latLong)
		if err != nil {
			log.Fatalf("(\"input#search\").Click(): %v", err)
		}
		if err = page.Locator(".list-group-item-wrapper").First().WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}); err != nil {
			log.Fatalf(".list-group-item-wrapper: %v", err)
		}
		firstAutoComplete, err := page.Locator(".list-group-item-wrapper").First().Locator(".body.list-two-line-text__title").TextContent()
		if err != nil {
			log.Fatalf(".list-group-item-wrapper: %v", err)
		}
		firstAutoCompleteDesc, err := page.Locator(".list-group-item-wrapper").First().Locator(".body-s.list-two-line-text__bottom-left").TextContent()
		log.Printf("firstAutoComplete %s %s %s", latLong, firstAutoComplete, firstAutoCompleteDesc)
		rec = append(rec, firstAutoComplete)
		rec = append(rec, firstAutoCompleteDesc)

		writer.Write(rec)
		if err = page.Locator("#page-container  .header__middle-area button.is-right").Click(); err != nil {
			log.Fatalf("button.is-right Click %v", err)
		}
	}

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
