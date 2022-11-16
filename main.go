package main

import (
	"fmt"
	"net/http"

	"weatherapp/src/constants"
	"weatherapp/src/controller"
)

func main() {
	http.HandleFunc("/weather", controller.GetCityWeather)
	fmt.Printf("Starting application on port %s\n", constants.PortNumber)
	_ = http.ListenAndServe(constants.PortNumber, nil)
}
