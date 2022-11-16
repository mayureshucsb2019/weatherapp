This is a weather app to display the weather at a particular location.

Software Requirements:
* Go version 1.18.4

HOW TO RUN THE APP
* Export APIID into your command line terminal using `export APPID=5b17073729841d77935a48ff6472f9fb`
* cd to the path where main.go is located
* Check if you have go.mod and go.sum files present
* Install all go dependencies using `go install`
* To run using main file do: `go run main.go`
* (optional) To build binary do: `go build`
* Binary will be generated with name weather_app 
* Run this binary using `./weather_app` 

SOLUTION EXPLANATION:
* The controller function makes use of two utility function `ConvertCityToCoordinates` and `GetWeatherReportForCoordinates`
* `ConvertCityToCoordinates` takes in a string input of city name and queries for its geo coordinates
*  We make use of simple caching mechanism with ttl of 24 hours for hot coordinates. This will help reduce the number of API calls to openwearthermap's geo api. 
* `ConvertCityToCoordinates` first looks into cache using `GetCache` utlity function to check if city coordinates are cached. Else it queries for the geo coordinates and adds it to cache.
* This fetched coordinates are used in `GetWeatherReportForCoordinates` as input to fetch weather information. Here, we can use caching but I choose not to use as weather keeps fluctuating and we might want to show the most recent weather. Nevertheless, TTL of 5 minutes or so can helpful if a particular city receives large number of queries. 