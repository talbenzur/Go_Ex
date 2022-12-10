package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type weather struct {
	minTemp  float32
	maxTemp  float32
	humidity int
	wind     float32
	rain     int
	city     string
}

// func mainCal(city string, days int) {

// 	valAtlas := mainAtlas(city, days)
// 	valTime := mainTime(city, days)
// 	valOpen := mainOpen(city, days)
// 	fmt.Println(getMaxAndMinTemp(averageWeather(valAtlas, valTime, valOpen, days)))
// 	fmt.Println(willItRain(averageDaysFromSite(valAtlas, valTime, valOpen, days), days))
// }

func getDataFromMoldules(city string, days int) ([]weather, []weather, []weather) {
	valAtlas := mainAtlas(city, days)
	valTime := mainTime(city, days)
	valOpen := mainOpen(city, days)
	return valAtlas, valTime, valOpen

}
func willItRainClient(city string, days int) []int {
	valAtlas, valTime, valOpen := getDataFromMoldules(city, days)
	return willItRain(averageDaysFromSite(valAtlas, valTime, valOpen, days), days)

}
func nextRainDay(city string, days int) int {
	valAtlas, valTime, valOpen := getDataFromMoldules(city, days)
	rainPercentages := willItRain(averageDaysFromSite(valAtlas, valTime, valOpen, days), days)

	for i := 0; i < days; i++ {
		if rainPercentages[i] > 50 {
			return i
		}
	}
	return -1
}

func avgTempretureClient(city string, days int) [2]float32 {

	valAtlas, valTime, valOpen := getDataFromMoldules(city, days)
	min, max := getMaxAndMinTemp(averageWeather(valAtlas, valTime, valOpen, days))
	// avg := (min + max) / 2
	var res [2]float32
	res[0] = min
	res[1] = max

	return res
}
func minMaxTempTill(city string, days int) ([]float32, []float32) {
	valAtlas, valTime, valOpen := getDataFromMoldules(city, days)
	// var res[2] [] float32
	weatherDays := averageDaysFromSite(valAtlas, valTime, valOpen, days)
	var min []float32
	var max []float32

	for i := 0; i < days; i++ {

		min = append(min, weatherDays[i].minTemp)
		max = append(max, weatherDays[i].maxTemp)

	}

	return min, max

}
func todaySum(city string) weather {
	valAtlas, valTime, valOpen := getDataFromMoldules(city, 1)
	w := averageWeather(valAtlas, valTime, valOpen, 1)

	return w
}

//=======================

// =======================
func willItRain(averageDaysFromSite []weather, days int) []int {
	var rain []int
	for _, dailyWeather := range averageDaysFromSite {
		rain = append(rain, dailyWeather.rain) //here we should add %
	}
	return rain
}
func getMaxAndMinTemp(averageWeather weather) (float32, float32) {
	return averageWeather.maxTemp, averageWeather.minTemp
}
func averageWeather(valAtlas []weather, valTime []weather, valOpen []weather, days int) weather {
	var averageWeather weather
	averageDaysFromSite := averageDaysFromSite(valAtlas, valTime, valOpen, days)
	for _, dailyWeather := range averageDaysFromSite {
		setAllDaysWeatherVals(&averageWeather, dailyWeather)
	}
	averageCal(&averageWeather, days)

	return averageWeather
}

func averageDaysFromSite(valAtlas []weather, valTime []weather, valOpen []weather, days int) []weather {
	var averageDaysFromSite []weather
	for i := 0; i < days; i++ {
		weatherType := weather{0, 0, 0, 0, 0, ""}
		setWeatherValsFromSites(&weatherType, valAtlas, i)
		setWeatherValsFromSites(&weatherType, valTime, i)
		setWeatherValsFromSites(&weatherType, valOpen, i)
		averageCal(&weatherType, 3)
		averageDaysFromSite = append(averageDaysFromSite, weatherType)
	}

	return averageDaysFromSite
}

func averageCal(weatherType *weather, num int) {
	weatherType.maxTemp = weatherType.maxTemp / float32(num)
	weatherType.minTemp = weatherType.minTemp / float32(num)
	weatherType.humidity = weatherType.humidity / num
	weatherType.rain = weatherType.rain / num
	weatherType.wind = weatherType.wind / float32(num)

}

func setAllDaysWeatherVals(weatherType *weather, weatherArr weather) {
	weatherType.maxTemp += weatherArr.maxTemp
	weatherType.minTemp += weatherArr.minTemp
	weatherType.humidity += weatherArr.humidity
	weatherType.wind += weatherArr.wind
	weatherType.rain += weatherArr.rain

}

func setWeatherValsFromSites(weatherType *weather, weatherArr []weather, index int) {

	weatherType.maxTemp += weatherArr[index].maxTemp
	weatherType.minTemp += weatherArr[index].minTemp
	weatherType.rain += weatherArr[index].rain
	weatherType.humidity += weatherArr[index].humidity
	weatherType.wind += weatherArr[index].wind
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
