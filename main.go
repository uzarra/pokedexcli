package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/uzarra/pokedexcli/internal/pokecache"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config, cache *pokecache.Cache, location string) error
}

type PokemonResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
}

type Response struct {
	Previous string     `json:"previous"`
	Next     string     `json:"next"`
	Results  []Location `json:"results"`
}

type Location struct {
	Name string `json:"name"`
}

type Config struct {
	Prev string
	Next string
}

func main() {
	url := "https://pokeapi.co/api/v2/location-area"
	config := Config{
		Prev: "",
		Next: url,
	}
	pagesCache := pokecache.NewCache(5 * time.Second)

	supportedCommands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Show this help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Map of locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Map of locations back",
			callback:    commandMapB,
		},
		"explore": {
			name:        "mapb",
			description: "Exploring location area",
			callback:    commandExplore,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		text := scanner.Text()
		values := strings.Split(text, " ")
		value, ok := supportedCommands[values[0]]
		if ok {
			location := ""
			if len(values) > 1 {
				location = values[1]
			}
			err := value.callback(&config, pagesCache, location)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func commandExit(config *Config, cache *pokecache.Cache, location string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, cache *pokecache.Cache, location string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex\n")
	return nil
}

func commandMap(config *Config, cache *pokecache.Cache, location string) error {
	if config.Next == "" {
		fmt.Println("you're on the last page")
		return nil
	}
	body, hasEntry := cache.Get(config.Next)
	if !hasEntry {
		var err error
		body, err = callApi(config.Next)
		if err != nil {
			return err
		}
		cache.Add(config.Next, body)
	}
	var response Response
	err := json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	prev := response.Previous
	next := response.Next
	config.Prev = prev
	config.Next = next
	for _, result := range response.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandMapB(config *Config, cache *pokecache.Cache, location string) error {
	if config.Prev == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	body, hasEntry := cache.Get(config.Prev)
	if !hasEntry {
		var err error
		body, err = callApi(config.Prev)
		if err != nil {
			return err
		}
		cache.Add(config.Prev, body)
	}
	var response Response
	err := json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	prev := response.Previous
	next := response.Next
	config.Prev = prev
	config.Next = next
	for _, result := range response.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandExplore(config *Config, cache *pokecache.Cache, location string) error {
	if location == "" {
		fmt.Println("location argument is empty")
		return nil
	}
	fullUrl := "https://pokeapi.co/api/v2/location-area/" + location
	body, hasEntry := cache.Get(fullUrl)
	if !hasEntry {
		var err error
		body, err = callApi(fullUrl)
		if err != nil {
			return err
		}
		cache.Add(fullUrl, body)
	}
	var response Response
	err := json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	var pokemonResponse PokemonResponse
	err = json.Unmarshal(body, &pokemonResponse)
	if err != nil {
		return err
	}
	for _, pokemonEncounter := range pokemonResponse.PokemonEncounters {
		fmt.Println(pokemonEncounter.Pokemon.Name)
	}
	return nil
}

func callApi(fullUrl string) ([]byte, error) {
	resp, err := http.Get(fullUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, nil
	}
	return nil, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
}
