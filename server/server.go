package server

import (
	"fmt"
	"net/http"

	Pokemon "github.com/XanderHK/Pokegoapi/server/src/app"
)

//
func Start() {
	routes()
	http.ListenAndServe(":9990", nil)
}

//
func routes() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		enableCors(&rw)

		jsonStringResult := Pokemon.GetAllPokemonNames()
		fmt.Fprint(rw, jsonStringResult)
	})

	http.HandleFunc("/pokemon", func(rw http.ResponseWriter, r *http.Request) {
		enableCors(&rw)

		// I suck at regex so we use query params
		pokemonId := r.URL.Query()["id"]
		if pokemonId[0] == "" {
			fmt.Fprint(rw, "error")
			return
		}

		jsonStringResult := Pokemon.GetPokemonById(pokemonId)

		fmt.Fprint(rw, jsonStringResult)
	})
}

//
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
