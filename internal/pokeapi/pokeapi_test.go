package pokeapi

import (
	"testing"
	"time"

	"github.com/dsledge1/Pokedex/internal/pokecache"
)

func TestGetLocations(t *testing.T) {
	cache := pokecache.NewCache(5 * time.Minute)

	// Test fetching locations from the API
	url := APIEndpoint + "location-area?offset=0&limit=20"
	response, err := GetLocations(url, cache)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if response.Count == 0 {
		t.Errorf("Expected count to be greater than 0")
	}
	if len(response.Results) == 0 {
		t.Errorf("Expected results to be non-empty")
	}

	// Test fetching the same URL again to hit the cache
	cachedResponse, err := GetLocations(url, cache)
	if err != nil {
		t.Fatalf("Expected no error on cached fetch, got %v", err)
	}
	if cachedResponse.Count != response.Count {
		t.Errorf("Expected cached count to match original count")
	}
	if len(cachedResponse.Results) != len(response.Results) {
		t.Errorf("Expected cached results length to match original results length")
	}
}
