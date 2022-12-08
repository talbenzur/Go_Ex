//go get github.com/PuerkitoBio/goquery

package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

//	type weather struct {
//		minTepm  float32
//		maxTemp  float32
//		humidity int
//		wind     float32
//		rain     int
//		city     string
//	}
type selector struct {
	minTemp  string
	maxTemp  string
	humidity string
	wind     string
	rain     string
}

var arr []weather

func main() {

	fmt.Println("Enter Your City Name: ")
	var city string
	fmt.Scanln(&city)
	fmt.Println("Select the number of days you would like to see the weather forecast for: ")
	var days int
	fmt.Scanln(&days)

	webPage := "https://www.weather-atlas.com/en/israel/" + city + "-long-term-weather-forecast"

	selecorPrefix := "body > div.container.container_white > div:nth-child(5) > div > div:nth-child(3) > div.col-lg-8 > div > div > div:nth-child(2) > div:nth-child() > "

	/* ================================================================================== Temprature =====================================================================*/
	MaxTempSelector := selecorPrefix + "div.col-sm-6.d-flex.align-items-center.justify-content-center.px-0 > div.col-4.text-center.p-0 > ul > li.fs-2.text-danger"
	MinTempSelector := selecorPrefix + "div.col-sm-6.d-flex.align-items-center.justify-content-center.px-0 > div.col-4.text-center.p-0 > ul > li.fs-3.text-primary"

	/* ================================================================================== Humidity =====================================================================*/
	HumiditySelector := selecorPrefix + "div.col-sm-6.px-0.ps-sm-2 > div > div:nth-child(1) > ul > li:nth-child(2)"

	/* ================================================================================== Rain =====================================================================*/
	RainSelector := selecorPrefix + "div.col-sm-6.px-0.ps-sm-2 > div > div:nth-child(2) > ul > li:nth-child(1) > span"
	/* ================================================================================== Wind =====================================================================*/
	WindSelector := selecorPrefix + "div.col-sm-6.px-0.ps-sm-2 > div > div:nth-child(1) > ul > li:nth-child(1) > span"

	// days := 5

	atlasWeatherSummary(webPage,
		selector{MinTempSelector,
			MaxTempSelector,
			HumiditySelector, WindSelector,
			RainSelector},
		days, city)
	// scrapeValues(getURL(webPage), MaxTempSelector, days)

}

func atlasWeatherSummary(webURL string, selectors selector, days int, city string) []weather {

	valuesByFieldsMap := make(map[string][]float32)

	doc := getURL(webURL)
	fields := reflect.VisibleFields(reflect.TypeOf(selectors))
	for _, field := range fields {

		r := reflect.ValueOf(selectors)
		selectorAsString := r.FieldByName(field.Name)
		values := scrapeValues(doc, selectorAsString.String(), days)
		// fmt.Println(field.Name, values)

		valuesByFieldsMap[field.Name] = values
	}

	// fmt.Println(valuesByFieldsMap)
	return convertMapToWethers(valuesByFieldsMap, days, city)

}

func convertMapToWethers(valuesMap map[string][]float32, days int, city string) []weather {
	var whetherList []weather

	for i := 0; i < days; i++ {
		whetherList = append(whetherList, weather{valuesMap["minTemp"][i], valuesMap["maxTemp"][i], int(valuesMap["humidity"][i]), valuesMap["wind"][i], int(valuesMap["rain"][i]), city})
	}
	fmt.Println(whetherList)
	return whetherList
}

func scrapeValues(doc *goquery.Document, selector string, days int) []float32 {
	var arr []string

	n := days + 1
	for i := 1; i < n; i++ {
		ValueSelector := selector[0:143] + strconv.Itoa(i) + selector[143:]

		values := doc.Find(ValueSelector).Map(func(i int, sel *goquery.Selection) string {
			return fmt.Sprintf("%s", sel.Text())
		})
		if len(values) != 0 {
			arr = append(arr, values[0])
		} else {
			arr = append(arr, "0")
		}

	}
	var floatArr []float32
	sampleRegexp := regexp.MustCompile(`\d{1,2}`)

	for index := range arr {

		values, err := strconv.ParseFloat(sampleRegexp.FindString(arr[index]), 32)
		if err != nil {
			// fmt.Println("Error during conversion")
		}
		floatArr = append(floatArr, float32(values))

	}
	return floatArr
}

// func getURL(webURL string) *goquery.Document {
// 	resp, err := http.Get(webURL)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer resp.Body.Close()

// 	doc, err := goquery.NewDocumentFromReader(resp.Body)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if resp.StatusCode != 200 {
// 		log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
// 	}
// 	return doc
// }
