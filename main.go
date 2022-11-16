package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jellydator/ttlcache/v2"
	"github.com/mitchellh/mapstructure"
)

// Port at which app runs.
const portNumber = ":3000"

// Limit for number of cities with same names to be queried to openweathermap.org to get co-ordinates.
const queryLimit = 10

// This is a secret and must be fetched from the evironment variable.
var appId = os.Getenv("APPID")

// Implements simple cache to store city co-ordinates
var cacheCoordinates ttlcache.SimpleCache = ttlcache.NewCache()

func getCache(city string) (coordinate, error) {
	var coord coordinate
	val, err := cacheCoordinates.Get(city)
	if err != nil {
		log.Printf("key %s, err received: %+v\n", city, err)
		return coordinate{}, errors.New("not cached")
	}
	// decode empty interface into struct
	mapstructure.Decode(val, &coord)
	log.Printf("Cache received: %+v and decoded %+v\n", val, coord)
	return coord, nil
}

// Stores the geo coordinates in lat, lon format
type coordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// Response received for conversion of city name to geo co-ordinates
type cityQueryResponseModel struct {
	Name        string
	Local_names interface{}
	Lat         float64
	Lon         float64
	Country     string
	State       string
}

func ConvertCityToCoordinates(city string) (coordinate, error) {
	log.Printf("ConvertCityToCoordinates: city received: %s\n", city)
	// Create request URL by adding city to it and do get request
	var coord coordinate
	coord, err := getCache(city)
	// if cached value is found
	if err == nil {
		log.Printf("City %s was cached %+v\n", city, coord)
		return coord, nil
	}

	requestURL := fmt.Sprintf("https://api.openweathermap.org/geo/1.0/direct?q=%s&limit=%d&appid=%s", city, queryLimit, appId)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Println(err)
		return coordinate{}, err
	}

	// read request body and parse it to struct
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return coordinate{}, err
	}

	var arrayOfCities []cityQueryResponseModel
	err = json.Unmarshal([]byte(body), &arrayOfCities)
	if err != nil {
		log.Fatalln(err)
		return coordinate{}, err
	}

	if len(arrayOfCities) == 0 {
		return coordinate{}, errors.New("no city found")
	}
	coord = coordinate{arrayOfCities[0].Lat, arrayOfCities[0].Lon}
	// set cache
	cacheCoordinates.Set(city, coord)
	return coord, nil
}

func GetWeatherReportForCoordinates(cord coordinate) ([]byte, error) {
	log.Printf("Coordinate received is: %+v\n", cord)
	// Create request URL by adding city to it and do get request
	requestURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s", cord.Lat, cord.Lon, appId)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	// read request body and parse it to struct
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return body, nil
}

func GetCityWeather(w http.ResponseWriter, r *http.Request) {
	// setting a cache for hot queries for a day
	cacheCoordinates.SetTTL(time.Duration(3600 * 24 * time.Second))
	if r.URL.Path != "/weather" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	city := r.URL.Query().Get("city")
	if city == "" {
		fmt.Fprintln(w, "City name not entered with the optional parameter")
		return
	}
	output, err := ConvertCityToCoordinates(city)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	report, err := GetWeatherReportForCoordinates(output)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fmt.Fprintln(w, string(report))
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "Page")
	}
}

func main() {
	http.HandleFunc("/weather", GetCityWeather)
	fmt.Printf("Starting application on port %s\n", portNumber)
	_ = http.ListenAndServe(portNumber, nil)
}
