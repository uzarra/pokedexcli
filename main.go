package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		text := scanner.Text()
		cleanText := strings.Fields(strings.ToLower(text))
		fmt.Printf("Your command was: %s\n", cleanText[0])
	}
}
