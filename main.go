package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {

	// regex1 := ">Rain<.*PercentageValue\">(.*)%" //
	// fmt.Println(regex1)

	c := colly.NewCollector(
		colly.AllowedDomains("https://weather.com/?Goto=Redirected"),
	)
	// Find and print all links
	c.OnHTML("Precip", func(e *colly.HTMLElement) {
		links := e.ChildAttrs("PercentageValue", "data-testid")
		fmt.Println(links)
		rainChance := e.ChildText("#PercentageValue")
		fmt.Println("rainChance", rainChance)

	})
	c.Visit("https://weather.com/weather/today/l/1113fba38e09773bd1c20c1f70b48f045cd0d1bd62eaa2d37a47c5cc58d12b7b")

}
