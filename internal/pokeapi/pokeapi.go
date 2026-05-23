package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	APIEndpoint = "https://pokeapi.co/api/v2/"
)

type APIResponse struct {
	Count    int    `json:"count"`
	Next     string `json: "next"`
	Previous string `json: "previous"`
	Results  []struct {
		Name string `json: "name"`
		URL  string `json: "url"`
	} `json: "results"`
}

func GetLocations(url string) (APIResponse, error) {
	//locationURL := APIEndpoint + "location-area?offset=" + fmt.Sprint(offset) + "&limit=20"
	res, err := http.Get(url)
	if err != nil {
		return APIResponse{}, err
	}
	defer res.Body.Close()

	var locations APIResponse
	err = json.NewDecoder(res.Body).Decode(&locations)
	if err != nil {
		return APIResponse{}, err
	}
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return locations, nil

}
