package functions

import (
	"encoding/json"
	"regexp"

	PokemonTypes "github.com/XanderHK/Pokegoapi/server/src/app/types"
)

// Gets the evolution chain url from species
func GetPokemonEvolutionUrl(url string) string {
	responseData := GetRequest(url)
	var responseObject PokemonTypes.PokemonSpeciesResponse
	json.Unmarshal(responseData, &responseObject)
	return responseObject.EvoChain.Url
}

// Gets all the evolutions of a pokemon and returns them.
func GetPokemonEvolutionChain(url string) []string {
	responseData := GetRequest(url)
	var responseObject PokemonTypes.Chain
	json.Unmarshal(responseData, &responseObject)

	evolutions := []string{responseObject.Chain.Species.Name}
	evolutions = append(evolutions, WalkEvolutionChain(responseObject.Chain.EvolvesTo)...)
	return evolutions
}

// Recursive function that searches the structure for all evolutions of a specific pokemon
func WalkEvolutionChain(evolesTo []PokemonTypes.EvolvesTo) []string {
	var evolutions []string

	if len(evolesTo) > 0 {
		evolutions = append(evolutions, evolesTo[0].Species.Name)
		evolutions = append(evolutions, WalkEvolutionChain(evolesTo[0].EvolvesTo)...)
	}

	return evolutions
}

// Gets the default sprite from pokémon, this is used to show the different evolutions
func GetPokemonSprite(name string) string {
	responseData := GetRequest("https://pokeapi.co/api/v2/pokemon/" + name)

	var responseObject PokemonTypes.PokemonSingleResponse
	json.Unmarshal(responseData, &responseObject)

	return responseObject.Sprites.Front
}

// Uses the species URL of the pokemon to get the first english description it finds
func GetPokemonDesc(url string) string {
	responseData := GetRequest(url)

	var responseObject PokemonTypes.PokemonDescriptions
	json.Unmarshal(responseData, &responseObject)

	var firstEnglishDesc string
	for _, desc := range responseObject.Entries {
		if desc.Language.Name == "en" {
			re := regexp.MustCompile(`\r?\n|\f`)
			firstEnglishDesc = re.ReplaceAllString(desc.FlavorText, " ")
			break
		}
	}

	return firstEnglishDesc
}

// function that gets all pokemons and returns the length a.k.a. the amount of pokemon
func GetPokemonEntries() int {
	url := "https://pokeapi.co/api/v2/pokemon-species/?limit=20000"
	responseData := GetRequest(url)
	var responseObject PokemonTypes.ResponseAll
	json.Unmarshal(responseData, &responseObject)

	amountOfEntries := len(responseObject.Pokemon)
	return amountOfEntries
}

//
func GetEvolutionSprites(evolutions []string) []string {
	var evolutionSprites []string
	for _, evo := range evolutions {
		evolutionSprites = append(evolutionSprites, GetPokemonSprite(evo))
	}
	return evolutionSprites
}

//
func GetPokemonTypes(types []PokemonTypes.PokemonTypes) []string {
	var returnTypes []string
	for _, pokemonType := range types {
		returnTypes = append(returnTypes, pokemonType.Type.Name)
	}
	return returnTypes
}

//
func GetPokemonStats(stats []PokemonTypes.PokemonStatsResponse) []PokemonTypes.PokemonStat {
	var returnStats []PokemonTypes.PokemonStat
	for _, pokemonStat := range stats {
		returnStats = append(returnStats, PokemonTypes.PokemonStat{Name: pokemonStat.Stat.Name, Amount: pokemonStat.BaseStat})
	}
	return returnStats
}
