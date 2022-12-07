package main

import (
	"fmt"

	"github.com/gocolly/colly"
)



func main() {

    regex1:= ">Rain<.*PercentageValue\">(.*)%" //https://weather.com/weather/tenday/l/cca0801d9062db761eb0521ed6bf5549ed4c546211046d09b9ac7ad09a6c556d#detailIndex5

	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)
    // Find and print all links
    c.OnHTML(".mw-parser-output", func(e *colly.HTMLElement) {
        links := e.ChildAttrs("a", "href")
        fmt.Println(links)
    })
    c.Visit("https://en.wikipedia.org/wiki/Web_scraping")

}