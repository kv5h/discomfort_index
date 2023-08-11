package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type GeoLocation struct {
	Latitude  float64
	Longitude float64
}

type WeatherData struct {
	City        string
	Humidity    int
	Temperature float64
}

type DiscomfortIndex struct {
	City        string
	Feeling     string
	Humidity    int
	Index       float64
	Temperature float64
}

type ApiReturnGeo struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

type ApiReturnWeather struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
		AirQuality struct {
			Co           float64 `json:"co"`
			No2          float64 `json:"no2"`
			O3           float64 `json:"o3"`
			So2          float64 `json:"so2"`
			Pm25         float64 `json:"pm2_5"`
			Pm10         float64 `json:"pm10"`
			UsEpaIndex   int     `json:"us-epa-index"`
			GbDefraIndex int     `json:"gb-defra-index"`
		} `json:"air_quality"`
	} `json:"current"`
}

func EntryPoint(ipAddress string, apiKey string) DiscomfortIndex {
	return getDiscomfortIndex(getCurrentWeather(getGeoLocation(ipAddress), apiKey))
}

func getGeoLocation(ipAddress string) GeoLocation {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ipAddress)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var apiReturnGeo ApiReturnGeo
	errGeo := json.Unmarshal(body, &apiReturnGeo)
	if errGeo != nil {
		fmt.Println(errGeo.Error())
	}

	var geoLocation GeoLocation
	geoLocation.Latitude = apiReturnGeo.Lat
	geoLocation.Longitude = apiReturnGeo.Lon

	return geoLocation
}

func getCurrentWeather(geolocation GeoLocation, apiKey string) WeatherData {
	latitudeStr := strconv.FormatFloat(geolocation.Latitude, 'f', -1, 64)
	longitudeStr := strconv.FormatFloat(geolocation.Longitude, 'f', -1, 64)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s,%s&aqi=yes", apiKey, latitudeStr, longitudeStr)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var apiReturnWeather ApiReturnWeather
	errWed := json.Unmarshal(body, &apiReturnWeather)
	if errWed != nil {
		fmt.Println(errWed.Error())
	}

	var weatherData WeatherData
	weatherData.Temperature = apiReturnWeather.Current.TempC
	weatherData.Humidity = apiReturnWeather.Current.Humidity
	weatherData.City = apiReturnWeather.Location.Name

	return weatherData
}

func getDiscomfortIndex(weatherData WeatherData) DiscomfortIndex {
	var discomfortIndex DiscomfortIndex
	var feeling string

	t := weatherData.Temperature
	h := float64(weatherData.Humidity)
	index := float64(0.81)*t + float64(0.01)*h*(float64(0.99)*t-float64(14.3)) + float64(46.3)

	switch {
	case index <= 55:
		feeling = "Cold"
		break
	case index <= 60:
		feeling = "Cold a little"
		break
	case index <= 65:
		feeling = "No feeling"
		break
	case index <= 70:
		feeling = "Feels Good"
		break
	case index <= 75:
		feeling = "Not Hot"
		break
	case index <= 80:
		feeling = "Hot a little"
		break
	case index <= 85:
		feeling = "Hot"
		break
	default:
		feeling = "Too Hot"
		break
	}

	discomfortIndex.City = weatherData.City
	discomfortIndex.Feeling = feeling
	discomfortIndex.Humidity = weatherData.Humidity
	discomfortIndex.Index = index
	discomfortIndex.Temperature = weatherData.Temperature

	return discomfortIndex
}
