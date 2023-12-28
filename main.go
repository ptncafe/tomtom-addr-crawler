package main

import (
	"fmt"
	"github.com/ptncafe/tomtom-addr-crawler/module"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "crawl-tomtom-addr",
				Aliases: []string{"c"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "input",
						Aliases: []string{"i"},
						Value:   "data/data.csv",
						Usage:   "Input CSV file",
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						//Value:   "data/result.csv",
						Usage: "Ouput CSV file",
					},
					&cli.StringFlag{
						Name:    "domain-url",
						Aliases: []string{"d"},
						Value:   "https://plan.tomtom.com/en/?p=10.82734,106.66315,9.55z&q=10.76397248,106.6881186",
						Usage:   "Tomtom Domain",
					},
				},
				Usage: "Crawl the address from the server Tomtom by browser",
				Action: func(c *cli.Context) error {
					err := module.CrawlAddressFromTomtom(c.Context,
						c.String("domain-url"),
						c.String("input"),
						c.String("output"),
					)
					if err != nil {
						errSM := module.SendMessageTelegramRetry(0, fmt.Sprintf("[ERROR] crawl-tomtom-addr %s got error: %v", hostName, err), 5)
						if errSM != nil {
							log.Fatal(errSM)
						}
						log.Fatalf("CrawlAddressFromTomtom %v", err)
					}
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
