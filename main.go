//go get github.com/PuerkitoBio/goquery

package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {

    webPage := "https://www.timeanddate.com/weather/israel/jerusalem"
    resp, err := http.Get(webPage)

    if err != nil {
        log.Fatal(err)
    }

    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
    }

    doc, err := goquery.NewDocumentFromReader(resp.Body)

    if err != nil {
        log.Fatal(err)
    }

    //title := doc.Find("title").Text()
    //fmt.Println(title)

    humiditySelector:= "table.table--left.table--inner-borders-rows tbody tr:nth-child(6) td:nth-child(2)"
    
    words := doc.Find(humiditySelector).Map(func(i int, sel *goquery.Selection) string {
        return fmt.Sprintf("%d: %s", i+1, sel.Text())
    })
    myword:=strings.Split( words[0:1][0], " ")
    //fmt.Println(myword[1][:len(myword[1])-1])

    sampleRegexp := regexp.MustCompile(`\d{1,3}`)
    match := sampleRegexp.FindString(myword[1])
    fmt.Println(match)
}


/*
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
*/