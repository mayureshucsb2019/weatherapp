package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"weatherapp/src/constants"
	"weatherapp/src/models"

	"github.com/jellydator/ttlcache/v2"
	"github.com/mitchellh/mapstructure"
)

// Implements simple cache to store city co-ordinates
var CacheCoordinates ttlcache.SimpleCache = ttlcache.NewCache()

func GetCache(city string) (models.Coordinate, error) {
	var coord models.Coordinate
	val, err := CacheCoordinates.Get(city)
	if err != nil {
		log.Printf("key %s, err received: %+v\n", city, err)
		return models.Coordinate{}, errors.New("not cached")
	}
	// decode empty interface into struct
	mapstructure.Decode(val, &coord)
	log.Printf("Cache received: %+v and decoded %+v\n", val, coord)
	return coord, nil
}

func ConvertCityToCoordinates(city string) (models.Coordinate, error) {
	log.Printf("ConvertCityToCoordinates: city received: %s\n", city)
	// Create request URL by adding city to it and do get request
	var coord models.Coordinate
	coord, err := GetCache(city)
	// if cached value is found
	if err == nil {
		log.Printf("City %s was cached %+v\n", city, coord)
		return coord, nil
	}

	requestURL := fmt.Sprintf("https://api.openweathermap.org/geo/1.0/direct?q=%s&limit=%d&appid=%s", city, constants.QueryLimit, constants.AppId)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Println(err)
		return models.Coordinate{}, err
	}

	// read request body and parse it to struct
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return models.Coordinate{}, err
	}

	var arrayOfCities []models.CityQueryResponseModel
	err = json.Unmarshal([]byte(body), &arrayOfCities)
	if err != nil {
		log.Fatalln(err)
		return models.Coordinate{}, err
	}

	if len(arrayOfCities) == 0 {
		return models.Coordinate{}, errors.New("no city found")
	}
	coord = models.Coordinate{arrayOfCities[0].Lat, arrayOfCities[0].Lon}
	// set cache
	CacheCoordinates.Set(city, coord)
	return coord, nil
}

func GetWeatherReportForCoordinates(cord models.Coordinate) ([]byte, error) {
	log.Printf("Coordinate received is: %+v\n", cord)
	// Create request URL by adding city to it and do get request
	requestURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s", cord.Lat, cord.Lon, constants.AppId)
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

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}
