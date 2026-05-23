package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/dsledge1/Pokedex/internal/pokeapi"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	config := &config{
		Limit:    20,
		Offset:   0,
		Next:     pokeapi.APIEndpoint + "location-area?offset=0&limit=20",
		Previous: pokeapi.APIEndpoint + "location-area?offset=0&limit=20",
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
	}

	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		command, ok := supportedCommands[input[0]]
		if ok {
			command.callback(config)
		} else {
			fmt.Println("Unknown command")
		}

	}

}

func commandExit(config *config) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func help(config *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}

func commandMap(config *config) error {
	loc, err := pokeapi.GetLocations(config.Next)
	if err != nil {
		return err
	}
	config.Next = loc.Next
	config.Previous = loc.Previous
	return nil
}

func commandMapb(config *config) error {
	loc, err := pokeapi.GetLocations(config.Previous)
	if err != nil {
		return err
	}
	config.Next = loc.Next
	config.Previous = loc.Previous
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Next     string `json: "next"`
	Previous string `json: "previous"`
	Limit    int    `json: "limit"`
	Offset   int    `json: "offset"`
}
