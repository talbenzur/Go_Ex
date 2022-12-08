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

type weather struct {
	minTepm  int
	maxTemp  int
	humidity int
	wind     int
	rain     int
	city     string
}

func main() {

	webPage := "https://www.timeanddate.com/weather/israel/jerusalem"
	selector := "table.table--left.table--inner-borders-rows tbody tr:nth-child(6) td:nth-child(2)"
	regex := `\d{1,3}`
	days := 5

	willItRain(webPage, selector, regex, days)
	weatherSummary(webPage, selector, regex)
}

func getURL(webURL string) *goquery.Document {
	resp, err := http.Get(webURL)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
	}
	return doc
}

func willItRain(webURL string, selector string, regex string, days int) string {
	doc := getURL(webURL)

	words := doc.Find(selector).Map(func(i int, sel *goquery.Selection) string {
		return fmt.Sprintf("%d: %s", i+1, sel.Text())
	})
	myword := strings.Split(words[0:1][0], " ")

	sampleRegexp := regexp.MustCompile(regex)
	match := sampleRegexp.FindString(myword[1])
	return match
}

func getRain(doc *goquery.Document, days int) {
	var arr []string

	n := days + 1
	for i := 1; i < n; i++ {
		rainPercentage := fmt.Sprint("table tbody tr:nth-child(", i, ") td:nth-child(9)")

		rainChance := doc.Find(rainPercentage).Map(func(i int, sel *goquery.Selection) string {
			return sel.Text()
		})

		// res := strings.Split(words[0:1][0], " ")
		arr = append(arr, rainChance[0])
	}
	fmt.Println(arr)
}

func getHumidity(doc *goquery.Document, days int) {
	var arr []string

	n := days + 1
	for i := 1; i < n; i++ {
		rainPercentage := fmt.Sprint("table tbody tr:nth-child(", i, ") td:nth-child(8)")

		rainChance := doc.Find(rainPercentage).Map(func(i int, sel *goquery.Selection) string {
			return sel.Text()
		})

		// res := strings.Split(words[0:1][0], " ")
		arr = append(arr, rainChance[0])
	}
	fmt.Println(arr)
}

func getTemp(doc *goquery.Document, days int) {
	var arrMax []string
	var arrMin []string

	n := days + 1
	for i := 1; i < n; i++ {
		minMaxTemp := fmt.Sprint("table tbody tr:nth-child(", i, ") td:nth-child(3)")

		maxTemps := doc.Find(minMaxTemp).Map(func(i int, sel *goquery.Selection) string {
			return sel.Text()
		})

		sampleRegexp1 := regexp.MustCompile(`(\d{1,2})`)
		matchMax := sampleRegexp1.FindString(maxTemps[0])

		minTemps := doc.Find(minMaxTemp).Map(func(i int, sel *goquery.Selection) string {
			return sel.Text()
		})
		sampleRegexp2 := regexp.MustCompile(`(\d{1,2})`)
		matchMin := sampleRegexp2.FindString(minTemps[0])

		arrMax = append(arrMax, matchMax)
		arrMin = append(arrMin, matchMin)

	}
	fmt.Println("max: ", arrMax)
	fmt.Println("min: ", arrMin)

}

func weatherSummary(webURL string, selector string, reg string) weather {
	webURL = webURL + "/ext"
	doc := getURL(webURL)
	getRain(doc, 14)
	getHumidity(doc, 14)
	getTemp(doc, 14)
	return weather{2, 3, 5, 6, 7, "jer"}

	// selector = "#qlook p"
	// words3 := doc.Find(selector).Map(func(i int, sel *goquery.Selection) string {
	// 	return fmt.Sprintf("%d: %s", i+1, sel.Text())
	// })

	// fmt.Println("words3: ", words3)
}
