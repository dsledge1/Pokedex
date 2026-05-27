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

// TODO - Determine if a single APIResponse struct is sufficient for all API calls, or if different structs for different data would be cleaner
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

type Pokemon struct {
	URL            string `json:"url"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Id             int    `json:"id"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		Base_stat int `json:"base_stat`
	}
}

func CatchPokemon(url string, cache *pokecache.Cache) (Pokemon, error) {
	//TODO - Return here to add a check for if the pokemon is on the most recent "explore"
	cachedData, ok := cache.Get(url)
	if ok {
		var cachedResponse Pokemon
		err := json.Unmarshal(cachedData, &cachedResponse)
		if err != nil {
			return Pokemon{}, err
		}
		return cachedResponse, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()
	var pokemon Pokemon
	bod := json.NewDecoder(res.Body)
	err = bod.Decode(&pokemon)
	if err != nil {
		return Pokemon{}, err
	}
	cachingPokemon, err := json.Marshal(pokemon)
	if err != nil {
		return Pokemon{}, err
	}
	cache.Add(url, cachingPokemon)
	fmt.Println("TEST OUTPUT" + pokemon.Name)
	return pokemon, nil
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
	var pokemonList APIResponse
	bod := json.NewDecoder(res.Body)
	err = bod.Decode(&pokemonList)
	if err != nil {
		return APIResponse{}, err
	}

	fmt.Println("Caching data...")
	cachingPokemon, err := json.Marshal(pokemonList)
	if err != nil {
		return APIResponse{}, err
	}
	cache.Add(url, cachingPokemon) //Get Pokemon struct right?
	for _, poke := range pokemonList.Pokemon_encounters {
		fmt.Println("- " + poke.Pokemon.Name)
	}
	return pokemonList, nil
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
