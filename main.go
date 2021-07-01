package main

import (
	"github.com/XanderHK/Pokegoapi/server"
	"github.com/XanderHK/Pokegoapi/server/src/commands"
)

func main() {
	server.Start()
	go commands.Scanner()
}
