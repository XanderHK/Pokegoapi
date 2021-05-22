package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type ResponseAll struct {
    Pokemon []PokemonAll `json:"results"`
}

type PokemonSpeciesResponse struct {
    EvoChain struct {
        Url string `json:"url"`
    } `json:"evolution_chain"`
}

type PokemonEvoChainResponse struct {
    Chain struct{
        EvolvesTo []struct{
            EvolvesTo []struct{
                Species struct {
                    Name string `json:"name"`
                } `json:"species"`
            } `json:"evolves_to"`
            Species struct {
                Name string `json:"name"`
            } `json:"species"`
        } `json:"evolves_to"`
        Species struct {
            Name string `json:"name"`
        } `json:"species"`
    } `json:"chain"`
}


// type Test struct {
    
// }

// type Example struct {
//     EvolvesTo []struct{
//         *Example
//     }
//     Species struct {
//         Name string `json:"name"`
//     } `json:"species"`
// }

type PokemonStatsResponse struct {
    BaseStat int `json:"base_stat"`
    Stat struct {
        Name string `json:"name"`
    } `json:"stat"`
}

type PokemonSingleResponse struct {
    Id int         `json:"id"`
    Name string    `json:"name"`
    Height float64 `json:"height"`
    Weight float64 `json:"weight"`
    Sprites PokemonSprites `json:"sprites"`
    Types []PokemonTypes `json:"types"`
    Species PokemonSpecies `json:"species"`
    Stats []PokemonStatsResponse `json:"stats"`
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

type PokemonSingleResult struct {
    Id int         `json:"id"`
    Name string    `json:"name"`
    Height float64 `json:"height"`
    Weight float64 `json:"weight"`
    Sprites PokemonSprites `json:"sprites"`
    Types []string `json:"types"`
    Species PokemonSpecies `json:"species"`
    Description string `json:"description"`
    Evolutions []string `json:"evolutions"`
    EvolutionSprites []string `json:"evolutionSprites"`
    Stats []PokemonStat `json:"stats"`
}

type PokemonSprites struct {
    Front string `json:"front_default"`
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

type PokemonDescriptions struct {
    Entries []struct {
        FlavorText string `json:"flavor_text"`
        Language struct {
            Name string `json:"name"`
        } `json:"language"`
    } `json:"flavor_text_entries"`
}

type PokemonDescription struct {
    FlavorText string
}

type PokemonSpecies struct {
    Url string `json:"url"`
}

type PokemonStat struct {
    Name string `json:"name"`
    Amount int `json:"amount"`
}


func main() {
    router()
}


func router () {
    http.HandleFunc("/", getAllPokemonNames)

    http.HandleFunc("/pokemon", getPokemonById)

    http.ListenAndServe(":9990", nil)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getAllPokemonNames(w http.ResponseWriter, r *http.Request) {
    enableCors(&w)
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
    enableCors(&w)
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

    var stats []PokemonStat
    for _, pokemonStat := range responseObject.Stats {
        stats = append(stats, PokemonStat{Name: pokemonStat.Stat.Name, Amount: pokemonStat.BaseStat})
    }

    result, _ := json.Marshal(PokemonSingleResult{
        Id: responseObject.Id, 
        Name: responseObject.Name, 
        Weight: responseObject.Weight, 
        Height: responseObject.Height, 
        Sprites: responseObject.Sprites,
        Types: types,
        Species: responseObject.Species,
        Description: description,
        Evolutions: evolutions,
        EvolutionSprites: evolutionSprites,
        Stats: stats,
    })
    fmt.Fprint(w, string(result))
}

func getPokemonDesc(url string) string{
    response, err := http.Get(url)
      if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }

    var responseObject PokemonDescriptions
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

func getPokemonEvolutionChain(url string) []string {
    responseSpecies, err := http.Get(url)
    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    responseDataSpecies, err := ioutil.ReadAll(responseSpecies.Body)
    if err != nil {
        log.Fatal(err)
    }

    var responseObjectSpecies PokemonSpeciesResponse
    json.Unmarshal(responseDataSpecies, &responseObjectSpecies)

    responseEvoChain, err := http.Get(responseObjectSpecies.EvoChain.Url)
    
    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    responseDataEvoChain, err := ioutil.ReadAll(responseEvoChain.Body)
    if err != nil {
        log.Fatal(err)
    }   

    var responseObjectEvoChain PokemonEvoChainResponse
    json.Unmarshal(responseDataEvoChain, &responseObjectEvoChain)

    var pokemonEvoChainNames []string

    pokemonEvoChainNames = append(pokemonEvoChainNames, responseObjectEvoChain.Chain.Species.Name)

    if len(responseObjectEvoChain.Chain.EvolvesTo) > 0 {
        pokemonEvoChainNames = append(pokemonEvoChainNames, responseObjectEvoChain.Chain.EvolvesTo[0].Species.Name)
        if len(responseObjectEvoChain.Chain.EvolvesTo[0].EvolvesTo) > 0 {
            pokemonEvoChainNames = append(pokemonEvoChainNames, responseObjectEvoChain.Chain.EvolvesTo[0].EvolvesTo[0].Species.Name)
        }
    }
    
    return pokemonEvoChainNames
}

func getPokemonSprite(name string) string {
    response, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + name)

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

    return responseObject.Sprites.Front
}

//