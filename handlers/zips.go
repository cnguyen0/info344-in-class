package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/info344-in-class/zipsvr/models"
)

type CityHandler struct {
	PathPrefix string
	Index      models.ZipIndex
}

// pointer and access the city Index
func (ch *CityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// we want to support this url : /zips/city-name
	// we need to grab that last token of the url

	// slicing syntax works on string because strings are arrays of bytes
	cityName := r.URL.Path[len(ch.PathPrefix):]
	cityName = strings.ToLower(cityName)
	if len(cityName) == 0 {
		// do not use log.fatal!!!

		// parameter: the writer, a string, and the http status
		http.Error(w, "please provide a city name", http.StatusBadRequest)
		return
	}

	w.Header().Add(headerContentType, contentTypeJSON)
	w.Header().Add(accessControlAllowOrigin, "*")
	zips := ch.Index[cityName]

	// we know its a json file because the content type is a json
	json.NewEncoder(w).Encode(zips)
}
