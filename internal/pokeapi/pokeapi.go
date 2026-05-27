package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dsledge1/Pokedex/internal/pokecache"
)

const (
	APIEndpoint = "https://pokeapi.co/api/v2/"
)

type APIResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
	Pokemon_encounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func GetPokemon(url string, cache *pokecache.Cache) (APIResponse, error) {
	fmt.Println("Scanning area for Pokemon")
	cachedData, ok := cache.Get(url)
	if ok {
		var cachedResponse APIResponse
		err := json.Unmarshal(cachedData, &cachedResponse)
		if err != nil {
			return APIResponse{}, err
		}
		// TODO Print the names of the Pokemon found in the cached response
		for _, poke := range cachedResponse.Pokemon_encounters {
			fmt.Println("- " + poke.Pokemon.Name)
		}
		return cachedResponse, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return APIResponse{}, err
	}
	defer res.Body.Close()
	var pokemon APIResponse
	bod := json.NewDecoder(res.Body)

	err = bod.Decode(&pokemon)
	if err != nil {
		return APIResponse{}, err
	}

	fmt.Println("Caching data...")
	cachingPokemon, err := json.Marshal(pokemon)
	if err != nil {
		return APIResponse{}, err
	}
	cache.Add(url, cachingPokemon) //Get Pokemon struct right?
	for _, poke := range pokemon.Pokemon_encounters {
		fmt.Println("- " + poke.Pokemon.Name)
	}
	return pokemon, nil
}

func GetLocations(url string, cache *pokecache.Cache) (APIResponse, error) {
	fmt.Print("Checking cache for map data\n")
	cachedData, ok := cache.Get(url)
	fmt.Println("cache check complete")
	if ok {
		var cachedResponse APIResponse
		err := json.Unmarshal(cachedData, &cachedResponse)
		if err != nil {
			return APIResponse{}, err
		}
		for _, cachedResponse := range cachedResponse.Results {
			fmt.Println(cachedResponse.Name)
		}
		return cachedResponse, nil
	}
	fmt.Println("No cached data found, calling API")
	res, err := http.Get(url)
	if err != nil {
		return APIResponse{}, err
	}
	defer res.Body.Close()

	var locations APIResponse
	bod := json.NewDecoder(res.Body)

	err = bod.Decode(&locations)
	if err != nil {
		return APIResponse{}, err
	}
	fmt.Println("Caching data...")
	cachingLocations, err := json.Marshal(locations)
	if err != nil {
		return APIResponse{}, err
	}
	cache.Add(url, cachingLocations)
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return locations, nil

}
