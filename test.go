package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// A Chain is used to describe the structure of a evolution chain response from the pokeapi
type Chain struct {
	Chain struct {
		EvolvesTo []EvolvesTo `json:"evolves_to"`
		Species struct {
			Name string `json:"name"`
		} `json:"species"`
	} `json:"chain"`
}

// The EvolvesTo type is used to describe the structure of the evolves_to attribute from the pokeapi 
// The EvolvesTo has a EvolvesTo attribute which has the type []EvolvesTo which makes it a recursive struct
// This way we can recursively get information from it without knowing the depth
type EvolvesTo struct {
	EvolvesTo []EvolvesTo `json:"evolves_to"`
	Species struct {
		Name string `json:"name"`
	} `json:"species"`
}	

// The Evolutions type is used to describe the structure of what we are trying to create
// The Marshal will turn this into a JSON object e.g. {evolutions:["item1", "item2", "item3"]}
type Evolutions struct {
	Evolutions []string `json:"evolutions"`
}



// Main gets a response from the pokeapi
// It will unmarshal/unserialize the JSON object assigns the rersult to a var with the Chain struct
// Then using walk it will get every Pokémon from respObject
// It will then marshal/serialize the Evolutions struct back into a byte slice with the found evolutions in it
// This will then be turned into a string and printed to the console
func main() {
	response, _ := http.Get("https://pokeapi.co/api/v2/evolution-chain/1/")
	data, _ := ioutil.ReadAll(response.Body)

	var respObject Chain
	json.Unmarshal(data, &respObject)

	evolutions := []string{respObject.Chain.Species.Name}
	evolutions = append(evolutions, Walk(respObject.Chain.EvolvesTo)...)

	result, _ := json.Marshal(Evolutions{
		Evolutions: evolutions,
	})

	fmt.Println(string(result))
}


// The Walk function will recursively traverse a array with EvolvesTo objects with in it
// It will check if the current array is bigger to zero to prevent a out of bounds exception
// Then it will grab the species name and append it to the evolutions array and then it will call itself using the next child as value
// The result of this will be appeneded onto the evolutions array too
// The end result will be a 1d array containing all the names
func Walk(evolesTo []EvolvesTo) []string {
	var evolutions []string

	if len(evolesTo) > 0 {
		evolutions = append(evolutions, evolesTo[0].Species.Name)
		evolutions = append(evolutions, Walk(evolesTo[0].EvolvesTo)...)
	}

	return evolutions
}