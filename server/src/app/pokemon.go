package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	PokemonTypes "github.com/XanderHK/Pokegoapi/server/src/app/types"
)

//
func GetAllPokemonNames() string {
	responseData := httpRequest("https://pokeapi.co/api/v2/pokemon-species/?limit=20000")
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
	responseData := httpRequest("https://pokeapi.co/api/v2/pokemon/" + pokemonId[0])

	var responseObject PokemonTypes.PokemonSingleResponse
	json.Unmarshal(responseData, &responseObject)

	if responseObject.Species.Url != "" {
		description := getPokemonDesc(responseObject.Species.Url)
		evolutions := getPokemonEvolutionChain(responseObject.Species.Url)
		var evolutionSprites []string
		for _, evo := range evolutions {
			evolutionSprites = append(evolutionSprites, getPokemonSprite(evo))
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

//
func getPokemonDesc(url string) string {
	responseData := httpRequest(url)

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

//
func getPokemonSprite(name string) string {
	responseData := httpRequest("https://pokeapi.co/api/v2/pokemon/" + name)

	var responseObject PokemonTypes.PokemonSingleResponse
	json.Unmarshal(responseData, &responseObject)

	return responseObject.Sprites.Front
}

//
func getPokemonEvolutionUrl(url string) string {
	responseData := httpRequest(url)
	var responseObject PokemonTypes.PokemonSpeciesResponse
	json.Unmarshal(responseData, &responseObject)
	return responseObject.EvoChain.Url
}

//
func getPokemonEvolutionChain(url string) []string {
	responseData := httpRequest(getPokemonEvolutionUrl(url))
	var responseObject PokemonTypes.Chain
	json.Unmarshal(responseData, &responseObject)

	evolutions := []string{responseObject.Chain.Species.Name}
	evolutions = append(evolutions, WalkEvolutionChain(responseObject.Chain.EvolvesTo)...)
	return evolutions
}

//
func WalkEvolutionChain(evolesTo []PokemonTypes.EvolvesTo) []string {
	var evolutions []string

	if len(evolesTo) > 0 {
		evolutions = append(evolutions, evolesTo[0].Species.Name)
		evolutions = append(evolutions, WalkEvolutionChain(evolesTo[0].EvolvesTo)...)
	}

	return evolutions
}

func httpRequest(url string) []byte {
	response, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	return responseData
}
