package main

import (
	"fmt"
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

func mainCal(city string, days int) {

	valAtlas := mainAtlas(city, days)
	valTime := mainTime(city, days)
	valOpen := mainOpen(city, days)
	fmt.Println(getMaxAndMinTemp(averageWeather(valAtlas, valTime, valOpen, days)))
	fmt.Println(willItRain(averageDaysFromSite(valAtlas, valTime, valOpen, days), days))
}

func willItRain(averageDaysFromSite []weather, days int) []string {
	var rain []string
	for _, dailyWeather := range averageDaysFromSite {
		rain = append(rain, string(dailyWeather.rain)+"% ")
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
