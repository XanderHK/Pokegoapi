package importPokemon

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	PokemonTypes "github.com/XanderHK/Pokegoapi/server/src/app/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.TODO()

// init func to initialize the db connection
func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("pokemon_storage").Collection("pokemon_collection")
}

// function that can be called from the main.go that makes the initial call for storing all pokemon in the db
func Pokemon() {
	start := time.Now()
	amountOfEntries := getPokemonEntries()
	onePercent := math.Round(float64(amountOfEntries) / 100)
	progress := 0

	for i := 1; i < amountOfEntries; i++ {
		pokemon := parseSinglePokemon(i)
		collection.InsertOne(ctx, pokemon)

		if i%int(onePercent) == 0 {
			progress++
		}

		if i == amountOfEntries {
			progress++
		}

		fmt.Printf("\r%d%% completed", progress)
	}
	end := time.Since(start)

	fmt.Printf("\n Importing all Pokémons took: %s ", end)
}

// function that gets all pokemons and returns the length a.k.a. the amount of pokemon
func getPokemonEntries() int {
	url := "https://pokeapi.co/api/v2/pokemon-species/?limit=20000"
	responseData := httpRequest(url)
	var responseObject PokemonTypes.ResponseAll
	json.Unmarshal(responseData, &responseObject)

	amountOfEntries := len(responseObject.Pokemon)
	return amountOfEntries
}

// Gets a pokemon and makes subsequent function calls / http request to get the other necessary data.
// Then it turns it into a byte slice that is a BSON object that can be interpreted and stored in mongodb
func parseSinglePokemon(id int) PokemonTypes.PokemonSingleResultBson {
	responseData := httpRequest("https://pokeapi.co/api/v2/pokemon/" + strconv.Itoa(id))
	var responseObject PokemonTypes.PokemonSingleResponse
	json.Unmarshal(responseData, &responseObject)

	description := getPokemonDesc(responseObject.Species.Url)

	evolutionUrl := getPokemonEvolutionUrl(responseObject.Species.Url)
	evolutions := getPokemonEvolutionChain(evolutionUrl)
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

	result := PokemonTypes.PokemonSingleResultBson{
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
	}

	return result
}

// Gets the evolution chain url from species
func getPokemonEvolutionUrl(url string) string {
	responseData := httpRequest(url)
	var responseObject PokemonTypes.PokemonSpeciesResponse
	json.Unmarshal(responseData, &responseObject)
	return responseObject.EvoChain.Url
}

// Gets all the evolutions of a pokemon and returns them.
func getPokemonEvolutionChain(url string) []string {
	responseData := httpRequest(url)
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
func getPokemonSprite(name string) string {
	responseData := httpRequest("https://pokeapi.co/api/v2/pokemon/" + name)

	var responseObject PokemonTypes.PokemonSingleResponse
	json.Unmarshal(responseData, &responseObject)

	return responseObject.Sprites.Front
}

// Uses the species URL of the pokemon to get the first english description it finds
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

// little wrapper function for http.get to avoid code duplication
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