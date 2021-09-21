package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"

	"github.com/XanderHK/Pokegoapi/server/src/functions"
)

type Command struct {
	CommandName string
}

var funcMap = map[string]interface{}{
	"import": importPokemon,
}

// Creates a scanner and calls a function depending on input
func Scanner() {
	commands := getValidCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if strings.Contains(input, "pokegoapi") {
			parts := strings.Split(input, " ")
			args := parts[1:]
			for _, arg := range args {
				for _, command := range commands {
					fmt.Println(command)
					if arg == command {
						fmt.Println(fmt.Sprintf("Executing %v", command))
						go callDynamically(command.(*Command).CommandName)
						break
					}
				}

				if !functions.Contains(commands, arg) {
					fmt.Println("Invalid argument")
				}
			}
		}
	}
}

func getValidCommands() []interface{} {
	files, _ := ioutil.ReadDir("./server/src/commands")
	var commandNames []interface{}
	for _, f := range files {
		var indexToSplit int
		parts := strings.Split(f.Name(), "")
		for i, part := range parts {
			runes := []rune(part)
			if unicode.IsUpper(runes[0]) {
				indexToSplit = i
				break
			}
		}
		commandName := f.Name()[0:indexToSplit]

		commandNames = append(commandNames, commandName)
	}
	return commandNames
}

func callDynamically(name string, args ...interface{}) {
	funcMap[name].(func())()
	// switch name {
	// case "import":
	// 	funcMap[name].(func())()
	// case "name":
	// 	funcMap["name"].(func(string))(args[0].(string))
	// }

}

// Check what command is being called
