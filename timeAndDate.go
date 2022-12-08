//go get github.com/PuerkitoBio/goquery

package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func mainTime(city string, days int) []weather {

	var webPage string = "https://www.timeanddate.com/weather/israel/"

	// willItRain(webPage, city, days)
	return weatherSummary(webPage, city, days)
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

func weatherSummary(webURL string, city string, days int) []weather {
	webURL = webURL + city + "/ext"
	doc := getURL(webURL)
	wind := getWind(doc, days)
	minTemp, maxTemp := getTemp(doc, days)
	humidity := getHumidity(doc, days)
	rain := getRain(doc, days)
	var weatherArray []weather
	for i := 0; i < days; i++ {
		weatherArray = append(weatherArray, weather{minTemp[i], maxTemp[i], humidity[i], wind[i], rain[i], city})
	}
	return weatherArray
}
