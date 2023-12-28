package module

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"time"
)

const domain_url = "https://plan.tomtom.com/en/?p=10.82734,106.66315,9.55z&q=10.76397248,106.6881186"

func CrawlAddressFromTomtom(ctx context.Context, domainUrl string, inputFile string, outputFile string) error {

	defer func() {
		if r := recover(); r != nil {
			errSM := SendMessageTelegram(0, fmt.Sprintf("[ERROR] crawl-tomtom-addr got panic: %v", r))
			if errSM != nil {
				log.Fatal(errSM)
			}
		}
	}()
	if len(domainUrl) <= 0 {
		domainUrl = domain_url
	}
	if len(inputFile) <= 0 {
		inputFile = "data/data.csv"
	}
	if len(outputFile) <= 0 {
		outputFile = fmt.Sprintf("data/result_%s.csv", time.Now().Format("2006_01_02_15_04_05Z07_00"))
	}
	pw, err := playwright.Run()
	if err != nil {
		logrus.WithError(err).Error("playwright.Run %v", err)
		return err
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		logrus.WithError(err).Error("Chromium.Launch %v", err)
		return err
	}
	page, err := browser.NewPage()
	if err != nil {
		logrus.WithError(err).Error("browser.NewPage %v", err)
		return err
	}
	if _, err = page.Goto(domainUrl); err != nil {
		logrus.WithError(err).Error("page.Goto %s %v", domainUrl, err)
		return err
	}

	f, err := os.Open(inputFile)
	if err != nil {
		logrus.WithError(err).Error("os.Open %s %v", domainUrl, err)
		return err
	}
	fileResult, err := os.Create(outputFile)
	if err != nil {
		logrus.WithError(err).Error("os.Create %s %v", domainUrl, err)
		return err
	}
	// remember to close the file at the end of the program
	defer f.Close()
	defer fileResult.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	writer := csv.NewWriter(fileResult)
	defer writer.Flush()
	logrus.Debugf("CrawlAddressFromTomtom init and started %s", domainUrl)
	index := -1
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		index++
		if index <= 0 {
			continue
		}
		if err = crawlDetail(ctx, page, rec, writer); err != nil {
			logrus.WithError(err).Error("crawlDetail %v %v", rec, err)
			return err
		}
	}

	if err = browser.Close(); err != nil {
		return err
	}
	if err = pw.Stop(); err != nil {
		return err
	}
	return nil
}
func crawlDetail(ctx context.Context, page playwright.Page, rec []string, writer *csv.Writer) (err error) {

	if err = page.Locator("input#search").Click(); err != nil {
		logrus.WithError(err).Error("input#search.Click %v", err)
		return err
	}
	if err = page.Locator("input#search").Clear(); err != nil {
		logrus.WithError(err).Error("input#search.Clear %v", err)
		return err
	}
	latLong := fmt.Sprintf("%s,%s", rec[1], rec[2])
	err = page.Locator("input#search").Fill(latLong)
	if err != nil {
		logrus.WithError(err).Error("input#search.Fill %v", err)
		return err
	}
	if err = page.Locator(".list-group-item-wrapper").First().WaitFor(playwright.LocatorWaitForOptions{State: playwright.WaitForSelectorStateVisible}); err != nil {
		logrus.WithError(err).Error("input#search.Fill %v", err)
		return err
	}
	firstAutoComplete, err := page.Locator(".list-group-item-wrapper").First().Locator(".body.list-two-line-text__title").TextContent()
	if err != nil {
		logrus.WithError(err).Error(".body.list-two-line-text__title TextContent %v", err)
		return err
	}
	firstAutoCompleteDesc, err := page.Locator(".list-group-item-wrapper").First().Locator(".body-s.list-two-line-text__bottom-left").TextContent()
	if err != nil {
		logrus.WithError(err).Error(".body-s.list-two-line-text__bottom-left TextContent %v", err)
		return err
	}

	logrus.Printf("firstAutoComplete %s %s %s", latLong, firstAutoComplete, firstAutoCompleteDesc)
	rec = append(rec, firstAutoComplete)
	rec = append(rec, firstAutoCompleteDesc)

	writer.Write(rec)
	if err = page.Locator("#page-container  .header__middle-area button.is-right").Click(); err != nil {
		return err
	}
	return nil
}
