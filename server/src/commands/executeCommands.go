package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Creates a scanner and calls a function depending on input
func Scanner() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if strings.Contains(input, "pokegoapi") {
			parts := strings.Split(input, " ")
			args := parts[1:]
			for _, arg := range args {
				switch arg {
				case "import":
					fmt.Println("Executing import")
					go ImportCommand()
				default:
					fmt.Println("Unknown argument")
				}
			}
		}
	}
}
