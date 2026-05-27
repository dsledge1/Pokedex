package main

import (
	"bufio"
	"fmt"
	"math/rand"
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
		fmt.Println("Name: " + pokemon + "\nHeight: ")
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
	fmt.Println(pok)
	exp := pok.BaseExperience
	fmt.Print("base experience: " + fmt.Sprint(exp) + "\n")
	fmt.Println("Pokemon is " + pok.Name + " with id " + fmt.Sprint(pok.Id))
	captureDifficulty := float64(exp) / float64(306)
	fmt.Printf("Capture difficulty: %.2f\n", captureDifficulty)
	random := rand.Float64()
	fmt.Printf("Random number: %.2f\n", random)
	if random >= captureDifficulty {

		fmt.Println(pokemon + " was caught!")
		config.Pokedex[pokemon] = pokeapi.Pokemon{
			URL:  pok.URL,
			Name: pok.Name,
			Id:   pok.Id,
		}
	} else {
		fmt.Println(pokemon + " escaped!")
		return nil
	}
	return nil
}

func commandExit(args []string, config *config) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

// TODO - Add additional help information for each command, and add a command to list all commands with descriptions
func help(args []string, config *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}

func commandMap(args []string, config *config) error {
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
