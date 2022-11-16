package controller

import (
	"fmt"
	"net/http"
	"time"

	"weatherapp/src/utils"
)

func GetCityWeather(w http.ResponseWriter, r *http.Request) {
	// setting a cache for hot queries for a day
	utils.CacheCoordinates.SetTTL(time.Duration(3600 * 24 * time.Second))
	if r.URL.Path != "/weather" {
		utils.ErrorHandler(w, r, http.StatusNotFound)
		return
	}
	city := r.URL.Query().Get("city")
	if city == "" {
		fmt.Fprintln(w, "City name not entered with the optional parameter")
		return
	}
	output, err := utils.ConvertCityToCoordinates(city)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	report, err := utils.GetWeatherReportForCoordinates(output)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fmt.Fprintln(w, string(report))
}
