package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"time"

	"github.com/dsledge1/Pokedex/internal/pokeapi"
	"github.com/dsledge1/Pokedex/internal/pokecache"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	config := &config{
		Limit:    20,
		Offset:   0,
		Next:     pokeapi.APIEndpoint + "location-area?offset=0&limit=20",
		Current:  pokeapi.APIEndpoint + "location-area?offset=0&limit=20",
		Previous: pokeapi.APIEndpoint + "location-area?offset=0&limit=20",
		Cache:    pokecache.NewCache(5 * time.Minute),
		Pokedex:  make(map[string]pokeapi.Pokemon),
	}

	supportedCommands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    help,
		},
		"map": {
			name:        "map",
			description: "Displays name of next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays name of previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays a list of all Pokemon located at a particular location",
			callback:    exploreLocation,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch the named Pokemon",
			callback:    catch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects the Pokedex entry for the named Pokemon",
			callback:    inspectPokedex,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all of the pokemon currently in your Pokedex",
			callback:    listPokedex,
		},
	}

	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}
		command, ok := supportedCommands[input[0]]
		if ok {
			err := command.callback(input, config)
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
		}

	}

}

func inspectPokedex(args []string, config *config) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: catch <pokemon-name>")
	}
	pokemon := args[1]
	if _, ok := config.Pokedex[pokemon]; ok {
		fmt.Println("Name: " + pokemon)
		fmt.Println("Height: " + fmt.Sprint(config.Pokedex[pokemon].Height))
		fmt.Println("Weight: " + fmt.Sprint(config.Pokedex[pokemon].Weight))
		fmt.Println("Stats:")
		for _, stat := range config.Pokedex[pokemon].Stats {
			fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.Base_stat)
		}
		fmt.Println("Types:")
		for _, typ := range config.Pokedex[pokemon].Types {
			fmt.Println("- " + typ.Type.Name)
		}
	} else {
		fmt.Println("You do not have " + pokemon + " in your Pokedex. Catch it first to inspect it!")
	}
	return nil
}

func listPokedex(args []string, config *config) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: pokedex")
	}
	if len(config.Pokedex) == 0 {
		fmt.Println("You do not have any Pokemon in your pokedex. Catch some!")
	} else {
		fmt.Println("Your Pokedex:")
		for _, pok := range config.Pokedex {
			fmt.Println("- " + pok.Name)
		}
	}
	return nil
}

func catch(args []string, config *config) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: catch <pokemon-name>")
	}
	pokemon := args[1]
	fmt.Println("Throwing a Pokeball at " + pokemon + "...")
	pok, err := pokeapi.CatchPokemon(pokeapi.APIEndpoint+"pokemon/"+pokemon, config.Cache)
	if err != nil {
		return err
	}
	exp := pok.BaseExperience
	captureDifficulty := float64(exp) / float64(306)
	random := rand.Float64()
	if random >= captureDifficulty {

		fmt.Println(pokemon + " was caught!")
		config.Pokedex[pokemon] = pokeapi.Pokemon{
			URL:            pok.URL,
			Name:           pok.Name,
			Id:             pok.Id,
			Height:         pok.Height,
			Weight:         pok.Weight,
			BaseExperience: pok.BaseExperience,
			Stats:          pok.Stats,
			Types:          pok.Types,
		}
	} else {
		fmt.Println(pokemon + " escaped!")
		return nil
	}
	return nil
}

func commandExit(args []string, config *config) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: exit")
	}
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

// TODO - Add additional help information for each command, and add a command to list all commands with descriptions
func help(args []string, config *config) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: help")
	}
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nmap <location>: Displays the next 20 locations\nmapb <location>: Displays the previous 20 locations\n\ncatch <pokemon-name>: Attempts to catch a wild Pokemon\n\ninspect <pokemon-name>: Inspects the named Pokemon if you have it in your Pokedex\npokedex: Lists the contents of your Pokedex\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}

func commandMap(args []string, config *config) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: map <location-area-name>")
	}
	loc, err := pokeapi.GetLocations(config.Next, config.Cache)
	if err != nil {
		return err
	}

	config.Current = config.Next
	config.Next = loc.Next
	config.Previous = loc.Previous
	return nil
}

func commandMapb(args []string, config *config) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: mapb <location-area-name>")
	}
	parsedURL, err := url.Parse(config.Current)
	queryParams := parsedURL.Query()
	if len(queryParams["offset"]) == 0 || queryParams["offset"][0] == "0" {
		return fmt.Errorf("Beginning of list - you cannot go back any further.")
	}
	loc, err := pokeapi.GetLocations(config.Previous, config.Cache)

	if err != nil {
		return err
	}
	config.Current = config.Previous
	config.Next = loc.Next
	config.Previous = loc.Previous
	return nil
}

func exploreLocation(args []string, config *config) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: explore <location-area-name>")
	}

	location := args[1]
	_, err := pokeapi.GetPokemon(pokeapi.APIEndpoint+"location-area/"+location, config.Cache)
	if err != nil {
		return err
	}
	return nil

}

type cliCommand struct {
	name        string
	description string
	callback    func([]string, *config) error
}

type config struct {
	Next     string `json: "next"`
	Current  string `json: current`
	Previous string `json: "previous"`
	Limit    int    `json: "limit"`
	Offset   int    `json: "offset"`
	Cache    *pokecache.Cache
	Pokedex  map[string]pokeapi.Pokemon
}

type Stats []struct {
	HP             int `json:"hp"`
	Attack         int `json:"attack"`
	Defense        int `json:"defense"`
	SpecialAttack  int `json:"special_attack"`
	SpecialDefense int `json:"special_defense"`
	Speed          int `json:"speed"`
}
type Type struct {
	Type string `json:"type"`
}
