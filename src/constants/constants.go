package constants

import "os"

// Port at which app runs.
const PortNumber = ":3000"

// Limit for number of cities with same names to be queried to openweathermap.org to get co-ordinates.
const QueryLimit = 10

// This is a secret and must be fetched from the evironment variable.
var AppId = os.Getenv("APPID")
