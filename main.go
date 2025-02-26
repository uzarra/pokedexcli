package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config) error
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
	Prev *string
	Next *string
}

func main() {
	url := "https://pokeapi.co/api/v2/location-area"
	config := Config{
		Prev: nil,
		Next: &url,
	}
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
			err := value.callback(&config)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func commandExit(config *Config) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex\n")
	return nil
}

func commandMap(config *Config) error {
	if config.Next == nil {
		fmt.Println("you're on the last page")
		return nil
	}
	resp, err := http.Get(*config.Next)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	prev := response.Previous
	next := response.Next
	config.Prev = &prev
	config.Next = &next
	for _, result := range response.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandMapB(config *Config) error {
	if config.Prev == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	resp, err := http.Get(*config.Prev)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	prev := response.Previous
	next := response.Next
	config.Prev = &prev
	config.Next = &next
	for _, result := range response.Results {
		fmt.Println(result.Name)
	}
	return nil
}
