package commands

import importPokemon "github.com/XanderHK/Pokegoapi/server/src/app/import"

// Invokes the Pokemon function from the importPokemon package
func ImportCommand() { importPokemon.Pokemon() }
