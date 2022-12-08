package main

import (
	"fmt"
	"log"
	"net/http"
	// "regexp"
	// "strings"
	// "github.com/PuerkitoBio/goquery"
	"encoding/json"
	"io/ioutil"
	"os"
)

type Response struct {
	Cod     string `json:"cod"`
	Message int    `json:"message"`
	Cnt     int    `json:"cnt"`
	List    []struct {
		Dt   int `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			SeaLevel  int     `json:"sea_level"`
			GrndLevel int     `json:"grnd_level"`
			Humidity  int     `json:"humidity"`
			TempKf    float64 `json:"temp_kf"`
		} `json:"main"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
			Gust  float64 `json:"gust"`
		} `json:"wind"`
		Visibility int `json:"visibility"`
		Pop        int `json:"pop"`
		Sys        struct {
			Pod string `json:"pod"`
		} `json:"sys"`
		DtTxt string `json:"dt_txt"`
	} `json:"list"`
	City struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country    string `json:"country"`
		Population int    `json:"population"`
		Timezone   int    `json:"timezone"`
		Sunrise    int    `json:"sunrise"`
		Sunset     int    `json:"sunset"`
	} `json:"city"`
}

type ResponseOne struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

func mainOpen() {

	lon, lat := getLonLat("jerusalem")

	lonStr := fmt.Sprintf("%f", lon)
	latStr := fmt.Sprintf("%f", lat)

	url := "https://api.openweathermap.org/data/2.5/forecast?lat=" + latStr + "&lon=" + lonStr + "&appid=104207efce9188324c22416dc6475c94&units=metric"
	response, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(string(responseData))

	var responseObject Response

	json.Unmarshal(responseData, &responseObject)
	dailySum(responseObject, 0)
	w := convertJsonToWeather(responseObject)
	// fmt.Println(w)
	fmt.Println(forcastRain(w, 4), getMaxTemp(w, 4), getMinTemp(w, 4))

}


func getLonLat(city string) (float64, float64) {
	url := "http://api.openweathermap.org/data/2.5/weather?q=" + city + ",israel&APPID=104207efce9188324c22416dc6475c94"
	responseOne, errOne := http.Get(url)

	if errOne != nil {
		fmt.Print(errOne.Error())
		os.Exit(1)
	}
	responseDataOne, err := ioutil.ReadAll(responseOne.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObjectOne ResponseOne

	json.Unmarshal(responseDataOne, &responseObjectOne)

	return responseObjectOne.Coord.Lon, responseObjectOne.Coord.Lat

}
func dailySum(responseObject Response, i int) weather {

	var w weather
	w.minTepm = float32(responseObject.List[i].Main.TempMin)
	w.maxTemp = float32(responseObject.List[i].Main.TempMax)
	w.wind = float32(responseObject.List[i].Wind.Speed)
	w.humidity = responseObject.List[i].Main.Humidity
	w.city = responseObject.City.Name

	return w

}

func convertJsonToWeather(responseObject Response) []weather {

	var weatherArr []weather
	for i := 0; i < len(responseObject.List); i = i + 8 { //here we should avg of day
		weatherArr = append(weatherArr, dailySum(responseObject, i))

	}

	return weatherArr

}

func forcastRain(wArr []weather, day int) []int {
	var rainRes []int

	for i := 0; i < day; i++ {
		rainRes = append(rainRes, wArr[i].rain)
	}
	return rainRes
}

func getMinTemp(wArr []weather, day int) []float32 {
	var minTempArr []float32

	for i := 0; i < day; i++ {
		minTempArr = append(minTempArr, wArr[i].minTepm)
	}
	return minTempArr
}
func getMaxTemp(wArr []weather, day int) []float32 {
	var maxTempArr []float32

	for i := 0; i < day; i++ {
		maxTempArr = append(maxTempArr, wArr[i].minTepm)
	}
	return maxTempArr
}
func avgTemp(wArr []weather, day int) float32 {
	var minRes float32
	var maxRes float32

	for i := 0; i < day; i++ {
		minRes += wArr[i].minTepm
		maxRes += wArr[i].maxTemp
	}
	avgMin := minRes / float32(day)
	avgMax := maxRes / float32(day)
	tempAvg := (avgMax + avgMin) / 2
	return tempAvg
}
func getMinMaxTemp(wArr []weather, day int) ([]float32, []float32) { //we should retrun these two arrays

	return getMinTemp(wArr, day), getMaxTemp(wArr, day)
}
