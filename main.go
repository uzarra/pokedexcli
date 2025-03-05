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
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config, cache *pokecache.Cache) error
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
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		text := scanner.Text()
		value, ok := supportedCommands[text]
		if ok {
			err := value.callback(&config, pagesCache)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func commandExit(config *Config, cache *pokecache.Cache) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, cache *pokecache.Cache) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex\n")
	return nil
}

func commandMap(config *Config, cache *pokecache.Cache) error {
	if config.Next == "" {
		fmt.Println("you're on the last page")
		return nil
	}
	body, hasEntry := cache.Get(config.Next)
	if !hasEntry {
		resp, err := http.Get(config.Next)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
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

func commandMapB(config *Config, cache *pokecache.Cache) error {
	if config.Prev == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	body, hasEntry := cache.Get(config.Prev)
	if !hasEntry {
		resp, err := http.Get(config.Prev)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
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
