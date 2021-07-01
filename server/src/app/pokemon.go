package app

import (
	"encoding/json"
	"strings"

	"github.com/XanderHK/Pokegoapi/server/src/app/functions"
	PokemonTypes "github.com/XanderHK/Pokegoapi/server/src/app/types"
)

//
func GetAllPokemonNames() string {
	responseData := functions.GetRequest("https://pokeapi.co/api/v2/pokemon-species/?limit=20000")
	var responseObject PokemonTypes.ResponseAll
	json.Unmarshal(responseData, &responseObject)

	var pokemonNamesAndIds []PokemonTypes.PokemonNameAndId
	for _, pokemon := range responseObject.Pokemon {

		var urlParts []string
		for _, v := range strings.Split(pokemon.Url, "/") {
			if v != "" {
				urlParts = append(urlParts, v)
			}
		}
		pokemonNamesAndIds = append(pokemonNamesAndIds, PokemonTypes.PokemonNameAndId{Name: pokemon.Name, Id: urlParts[len(urlParts)-1]})
	}

	result, _ := json.Marshal(PokemonTypes.PokemonNamesAndIds{Pokemon: pokemonNamesAndIds})
	return string(result)
}

//
func GetPokemonById(pokemonId []string) string {
	responseData := functions.GetRequest("https://pokeapi.co/api/v2/pokemon/" + pokemonId[0])

	var responseObject PokemonTypes.PokemonSingleResponse
	json.Unmarshal(responseData, &responseObject)

	if responseObject.Species.Url != "" {
		description := functions.GetPokemonDesc(responseObject.Species.Url)
		evolutions := functions.GetPokemonEvolutionChain(responseObject.Species.Url)
		var evolutionSprites []string
		for _, evo := range evolutions {
			evolutionSprites = append(evolutionSprites, functions.GetPokemonSprite(evo))
		}

		var types []string
		for _, pokemonType := range responseObject.Types {
			types = append(types, pokemonType.Type.Name)
		}

		var stats []PokemonTypes.PokemonStat
		for _, pokemonStat := range responseObject.Stats {
			stats = append(stats, PokemonTypes.PokemonStat{Name: pokemonStat.Stat.Name, Amount: pokemonStat.BaseStat})
		}

		result, _ := json.Marshal(PokemonTypes.PokemonSingleResult{
			Id:               responseObject.Id,
			Name:             responseObject.Name,
			Weight:           responseObject.Weight,
			Height:           responseObject.Height,
			Sprites:          responseObject.Sprites,
			Types:            types,
			Species:          responseObject.Species,
			Description:      description,
			Evolutions:       evolutions,
			EvolutionSprites: evolutionSprites,
			Stats:            stats,
		})

		return string(result)
	}

	return "oops something went wrong"
}
