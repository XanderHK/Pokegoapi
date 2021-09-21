package types

//
type ResponseAll struct {
	Pokemon []PokemonAll `json:"results"`
}

//
type PokemonSpeciesResponse struct {
	EvoChain struct {
		Url string `json:"url"`
	} `json:"evolution_chain"`
}

//
type PokemonStatsResponse struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

//
type PokemonSingleResponse struct {
	Id      int                    `json:"id"`
	Name    string                 `json:"name"`
	Height  float64                `json:"height"`
	Weight  float64                `json:"weight"`
	Sprites PokemonSprites         `json:"sprites"`
	Types   []PokemonTypes         `json:"types"`
	Species PokemonSpecies         `json:"species"`
	Stats   []PokemonStatsResponse `json:"stats"`
}

//
type PokemonAll struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

//
type PokemonNamesAndIds struct {
	Pokemon []PokemonNameAndId `json:"results"`
}

//
type PokemonNameAndId struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

//
type PokemonSingleResult struct {
	Id               int            `json:"id"`
	Name             string         `json:"name"`
	Height           float64        `json:"height"`
	Weight           float64        `json:"weight"`
	Sprites          PokemonSprites `json:"sprites"`
	Types            []string       `json:"types"`
	Species          PokemonSpecies `json:"species"`
	Description      string         `json:"description"`
	Evolutions       []string       `json:"evolutions"`
	EvolutionSprites []string       `json:"evolutionSprites"`
	Stats            []PokemonStat  `json:"stats"`
}

//
type PokemonSingleResultBson struct {
	Id               int            `bson:"id"`
	Name             string         `bson:"name"`
	Height           float64        `bson:"height"`
	Weight           float64        `bson:"weight"`
	Sprites          PokemonSprites `bson:"sprites"`
	Types            []string       `bson:"types"`
	Species          PokemonSpecies `bson:"species"`
	Description      string         `bson:"description"`
	Evolutions       []string       `bson:"evolutions"`
	EvolutionSprites []string       `bson:"evolutionSprites"`
	Stats            []PokemonStat  `bson:"stats"`
}

//
type PokemonSprites struct {
	Front string `json:"front_default"`
	Other struct {
		OfficialArtwork struct {
			FrontDefault string `json:"front_default"`
		} `json:"official-artwork"`
	} `json:"other"`
}

//
type PokemonTypes struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

//
type PokemonDescriptions struct {
	Entries []struct {
		FlavorText string `json:"flavor_text"`
		Language   struct {
			Name string `json:"name"`
		} `json:"language"`
	} `json:"flavor_text_entries"`
}

//
type PokemonDescription struct {
	FlavorText string
}

//
type PokemonSpecies struct {
	Url string `json:"url"`
}

//
type PokemonStat struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

// A Chain is used to describe the structure of a evolution chain response from the pokeapi
type Chain struct {
	Chain struct {
		EvolvesTo []EvolvesTo `json:"evolves_to"`
		Species   struct {
			Name string `json:"name"`
		} `json:"species"`
	} `json:"chain"`
}

// The EvolvesTo type is used to describe the structure of the evolves_to attribute from the pokeapi
// The EvolvesTo has a EvolvesTo attribute which has the type []EvolvesTo which makes it a recursive struct
// This way we can recursively get information from it without knowing the depth
type EvolvesTo struct {
	EvolvesTo []EvolvesTo `json:"evolves_to"`
	Species   struct {
		Name string `json:"name"`
	} `json:"species"`
}

// The Evolutions type is used to describe the structure of what we are trying to create
// The Marshal will turn this into a JSON object e.g. {evolutions:["item1", "item2", "item3"]}
type Evolutions struct {
	Evolutions []string `json:"evolutions"`
}
