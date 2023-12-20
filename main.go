package main

import (
	"log"
	"time"

	"github.com/playwright-community/playwright-go"
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
	err = page.Locator("input#search").Click()
	if err != nil {
		log.Fatalf("(\"input#search\").Click(): %v", err)
	}
	err = page.Locator("input#search").Clear()
	if err != nil {
		log.Fatalf("(\"input#search\").Click(): %v", err)
	}

	err = page.Locator("input#search").Fill("10.76131214,106.6888455")
	if err != nil {
		log.Fatalf("(\"input#search\").Click(): %v", err)
	}
	time.Sleep(1000)
	locatorList, err := page.Locator("#page-container > div.page-map > div.panels-container > div.panels-container__panel.panel-template.panel-template--top.panel-template--search > div.animation-wrapper.animation-wrapper--search-routes > div > div > div > div > div > div.simplebar-wrapper > div.simplebar-mask > div > div > div > div > div:nth-child(1)").All()
	log.Printf(".search-results.search-results--search .lists-item.lists-item--two-line %v", locatorList)

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
