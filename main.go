package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type weather struct {
	minTepm  float32
	maxTemp  float32
	humidity int
	wind     float32
	rain     int
	city     string
}

func main() {

	webPage := "https://www.timeanddate.com/weather/israel/"
	selector := "table.table--left.table--inner-borders-rows tbody tr:nth-child(6) td:nth-child(2)"
	regex := `\d{1,3}`
	days := 5

	willItRain(webPage, selector, regex, "jerusalem", days)
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

func willItRain(webURL string, selector string, regex string, city string, days int) string {
	doc := getURL(webURL + city)

	words := doc.Find(selector).Map(func(i int, sel *goquery.Selection) string {
		return fmt.Sprintf("%d: %s", i+1, sel.Text())
	})
	myword := strings.Split(words[0:1][0], " ")

	sampleRegexp := regexp.MustCompile(regex)
	match := sampleRegexp.FindString(myword[1])
	return match
}

func getRain(doc *goquery.Document, days int) []int {
	var arr []int

	n := days + 1
	for i := 1; i < n; i++ {
		rainPercentage := fmt.Sprint("table tbody tr:nth-child(", i, ") td:nth-child(9)")

		rainChance := doc.Find(rainPercentage).Map(func(i int, sel *goquery.Selection) string {
			return sel.Text()
		})
		sampleRegexp := regexp.MustCompile(`\d{1,2}`)
		match := sampleRegexp.FindString(rainChance[0])
		intVal, err := strconv.Atoi(match)

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		arr = append(arr, intVal)

	}
	return arr
}

func getHumidity(doc *goquery.Document, days int) []int {
	var arr []int

	n := days + 1
	for i := 1; i < n; i++ {
		rainPercentage := fmt.Sprint("table tbody tr:nth-child(", i, ") td:nth-child(8)")

		rainChance := doc.Find(rainPercentage).Map(func(i int, sel *goquery.Selection) string {
			return sel.Text()
		})
		sampleRegexp := regexp.MustCompile(`\d{1,2}`)
		match := sampleRegexp.FindString(rainChance[0])
		intVal, err := strconv.Atoi(match)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		arr = append(arr, intVal)

	}
	return arr
}

func getTemp(doc *goquery.Document, days int) ([]float32, []float32) {
	var arrMax []float32
	var arrMin []float32

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

		maxVal, err := strconv.ParseFloat(matchMax, 32)
		minVal, err := strconv.ParseFloat(matchMin[0][2], 32)

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		arrMax = append(arrMax, float32(maxVal))
		arrMin = append(arrMin, float32(minVal))

	}
	return arrMax, arrMin
}

func averageTemp(doc *goquery.Document, city string, days int) float32 {
	var averageDayTemp []float32
	arrMax, arrMin := getTemp(doc, days)
	var averageFinal float32 = 0

	for i := 0; i < len(arrMax); i++ {
		average := (arrMax[i] + arrMin[i]) / 2
		averageDayTemp = append(averageDayTemp, average)
	}
	for i := 0; i < len(averageDayTemp); i++ {
		averageFinal = averageDayTemp[i] + averageFinal
	}
	return (averageFinal / float32(days))
}

func getWind(doc *goquery.Document, days int) []float32 {
	var windArr []float32
	n := days + 1
	for i := 1; i < n; i++ {
		wind := fmt.Sprint("table tbody tr:nth-child(", i, ") td:nth-child(6)")

		windGet := doc.Find(wind).Map(func(i int, sel *goquery.Selection) string {
			return sel.Text()
		})
		sampleRegexp := regexp.MustCompile(`\d{1,2}`)
		match := sampleRegexp.FindString(windGet[0])
		floatVal, err := strconv.ParseFloat(match, 32)

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		windArr = append(windArr, float32(floatVal))

	}

	return windArr
}

func weatherSummary(webURL string, selector string, reg string) weather {
	webURL = webURL + "jerusalem" + "/ext"
	doc := getURL(webURL)
	wind := getWind(doc, 1)
	minTemp, maxTemp := getTemp(doc, 1)
	humidity := getHumidity(doc, 1)
	rain := getRain(doc, 1)
	var todayWeather = weather{minTemp[0], maxTemp[0], humidity[0], wind[0], rain[0], "jerusalem"}
	fmt.Println(todayWeather)
	return todayWeather
}
