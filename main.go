package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func main() {
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
			err := value.callback()
			if err != nil {
				fmt.Print(err.Error())
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex\n")
	return nil
}
