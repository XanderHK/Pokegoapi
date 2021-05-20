package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type ResponseAll struct {
    Pokemon []PokemonAll `json:"results"`
}

type PokemonAll struct {
    Name string `json:"name"`
    Url string  `json:"url"`
}

type PokemonNamesAndIds struct {
    Pokemon []PokemonNameAndId `json:"results"`
}

type PokemonNameAndId struct {
    Name string `json:"name"`
    Id string  `json:"id"`
}

type PokemonSingleResponse struct {
    Id int         `json:"id"`
    Name string    `json:"name"`
    Height float64 `json:"height"`
    Weight float64 `json:"weight"`
    Sprites PokemonSprites `json:"sprites"`
    Types []PokemonTypes `json:"types"`
    Species PokemonSpecies `json:"species"`
}

type PokemonSingleResult struct {
    Id int         `json:"id"`
    Name string    `json:"name"`
    Height float64 `json:"height"`
    Weight float64 `json:"weight"`
    Sprites PokemonSprites `json:"sprites"`
    Types []PokemonTypes `json:"types"`
    Species PokemonSpecies `json:"species"`
    Description PokemonDescription `json:"description"`
}

type PokemonSprites struct {
    Other struct {
         OfficialArtwork struct {
            FrontDefault string `json:"front_default"`
         } `json:"official-artwork"`
    } `json:"other"`
}

type PokemonTypes struct {
    Type struct {
        Name string `json:"name"`
    } `json:"type"`      
}

type PokemonDescription struct {
    Entries []struct {
        FlavorText string `json:"flavor_text"`
        Language struct {
            Name string `json:"name"`
        } `json:"language"`
    } `json:"flavor_text_entries"`
}

type PokemonSpecies struct {
    Url string `json:"url"`
}


func main() {
    router()
}


func router () {
    http.HandleFunc("/", getAllPokemonNames)

    http.HandleFunc("/pokemon", getPokemonById)

    http.ListenAndServe(":9990", nil)
}

func getAllPokemonNames(w http.ResponseWriter, r *http.Request) {
    response, err := http.Get("https://pokeapi.co/api/v2/pokemon-species/?limit=20000") 
    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }

    var responseObject ResponseAll

    json.Unmarshal(responseData, &responseObject)

    var pokemonNamesAndIds []PokemonNameAndId
    for _, pokemon := range responseObject.Pokemon {
    
        var urlParts []string
        for _, v := range strings.Split(pokemon.Url, "/") {
            if v != "" {
                urlParts = append(urlParts, v)
            }
        }
        pokemonNamesAndIds = append(pokemonNamesAndIds, PokemonNameAndId{Name: pokemon.Name, Id: urlParts[len(urlParts) - 1]})
    }

    result, _ := json.Marshal(PokemonNamesAndIds{Pokemon: pokemonNamesAndIds})
    fmt.Fprint(w, string(result))
}

func getPokemonById(w http.ResponseWriter, r *http.Request) {
    pokemonId := r.URL.Query()["id"]
    if pokemonId[0] == "" {
        fmt.Fprint(w, "error")
        return
    }

    response, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + pokemonId[0]) 
    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }

    var responseObject PokemonSingleResponse
    json.Unmarshal(responseData, &responseObject)

    descriptions := getPokemonDesc(responseObject.Species.Url)

    result, _ := json.Marshal(PokemonSingleResult{
        Id: responseObject.Id, 
        Name: responseObject.Name, 
        Weight: responseObject.Weight, 
        Height: responseObject.Height, 
        Sprites: responseObject.Sprites,
        Types: responseObject.Types,
        Species: responseObject.Species,
        Description: descriptions,
    })
    fmt.Fprint(w, string(result))
}

func getPokemonDesc(url string) PokemonDescription{
    response, err := http.Get(url)
      if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }

    var responseObject PokemonDescription
    json.Unmarshal(responseData, &responseObject)

    return responseObject
}
