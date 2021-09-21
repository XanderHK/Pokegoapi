package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/XanderHK/Pokegoapi/environment"
	"github.com/XanderHK/Pokegoapi/server/src/functions"
	PokemonTypes "github.com/XanderHK/Pokegoapi/server/src/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.TODO()

// init func to initialize the db connection
func init() {
	clientOptions := options.Client().
		ApplyURI(environment.GetEnvVariable("DATABASE_HOST"))

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database(environment.GetEnvVariable("DATABASE_NAME")).
		Collection(environment.GetEnvVariable("DATABASE_COLLECTION"))
}

// function that can be called from the main.go that makes the initial call for storing all pokemon in the db
func importPokemon() {
	start := time.Now()
	amountOfEntries := functions.GetPokemonEntries()
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

// Gets a pokemon and makes subsequent function calls / http request to get the other necessary data.
// Then it turns it into a byte slice that is a BSON object that can be interpreted and stored in mongodb
func parseSinglePokemon(id int) PokemonTypes.PokemonSingleResultBson {
	responseData := functions.GetRequest("https://pokeapi.co/api/v2/pokemon/" + strconv.Itoa(id))
	var responseObject PokemonTypes.PokemonSingleResponse
	json.Unmarshal(responseData, &responseObject)

	description := functions.GetPokemonDesc(responseObject.Species.Url)
	evolutionUrl := functions.GetPokemonEvolutionUrl(responseObject.Species.Url)
	evolutions := functions.GetPokemonEvolutionChain(evolutionUrl)
	evolutionSprites := functions.GetEvolutionSprites(evolutions)
	types := functions.GetPokemonTypes(responseObject.Types)
	stats := functions.GetPokemonStats(responseObject.Stats)

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
