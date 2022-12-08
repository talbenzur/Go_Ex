//go get github.com/PuerkitoBio/goquery

package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type weather struct {
	temp     tempPerDay
	humidity string
	wind     string
	rain     string
}

type tempPerDay struct {
	minTepm int
	maxTemp int
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

func getRain(doc *goquery.Document, days int) []string {
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
	return arr
}

func getHumidity(doc *goquery.Document, days int) []string {
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
	return arr
}

func getTemp(doc *goquery.Document, days int) ([]string, []string) {
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
		sampleRegexp2 := regexp.MustCompile(`(\d{1,2} \/ )(\d{1,2})`)
		matchMin := sampleRegexp2.FindAllStringSubmatch(minTemps[0], -1)

		arrMax = append(arrMax, matchMax)
		arrMin = append(arrMin, matchMin[0][2])

	}
	return arrMax, arrMin
}

func averageTemp(doc *goquery.Document, city string, days int) int {
	var arrMax []string
	var arrMin []string
	var averageDayTemp []int
	arrMax, arrMin = getTemp(doc, days)
	var averageFinal int = 0

	for i := 0; i < len(arrMax); i++ {
		intVar1, err1 := strconv.Atoi(arrMax[i])
		intVar2, err2 := strconv.Atoi(arrMin[i])
		if err1 == nil && err2 == nil {
			average := (intVar1 + intVar2) / 2
			averageDayTemp = append(averageDayTemp, average)
		}
	}

	for i := 0; i < len(averageDayTemp); i++ {
		averageFinal = averageDayTemp[i] + averageFinal
	}
	return (averageFinal / days)
}

func tempArray(doc *goquery.Document, city string, days int) []tempPerDay {
	var arrMax []string
	var arrMin []string
	var tempPerDayArr []tempPerDay

	arrMax, arrMin = getTemp(doc, days)

	for i := 0; i < days; i++ {
		intVar1, err1 := strconv.Atoi(arrMax[i])
		intVar2, err2 := strconv.Atoi(arrMin[i])
		if err1 == nil && err2 == nil {
			tempPerDayArr = append(tempPerDayArr, tempPerDay{intVar1, intVar2})
		}
	}
	return tempPerDayArr
}

func getWind(doc *goquery.Document, days int) []string {
	var windArr []string
	n := days + 1
	for i := 1; i < n; i++ {
		wind := fmt.Sprint("table tbody tr:nth-child(", i, ") td:nth-child(6)")

		windGet := doc.Find(wind).Map(func(i int, sel *goquery.Selection) string {
			return sel.Text()
		})
		windArr = append(windArr, windGet[0])
	}
	return windArr
}

func weatherSummary(webURL string, selector string, reg string) weather {
	webURL = webURL + "/ext"
	doc := getURL(webURL)
	wind := getWind(doc, 1)
	temp := tempArray(doc, "jerusalem", 1)
	humidity := getHumidity(doc, 1)
	rain := getRain(doc, 1)
	var todayWeather = weather{temp[0], humidity[0], wind[0], rain[0]}
	fmt.Println(todayWeather)
	return todayWeather
}
