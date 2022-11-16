package controller

import (
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
		http.Error(w, "Enter city as optional parameter.", http.StatusUnprocessableEntity)
		return
	}
	coord, err := utils.ConvertCityToCoordinates(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	report, err := utils.GetWeatherReportForCoordinates(coord)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(report)
}
