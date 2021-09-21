package main

import (
	"fmt"

	"github.com/XanderHK/Pokegoapi/environment"
	"github.com/XanderHK/Pokegoapi/server"
	"github.com/XanderHK/Pokegoapi/server/src/commands"
)

func main() {
	result := environment.GetEnvVariable("SERVER_PORT")
	fmt.Println(result)
	// var times []string

	// for i := 1; i <= 898; i++ {
	// 	start := time.Now()

	// 	resp, _ := http.Get("https://pokeapi.co/api/v2/pokemon/" + strconv.Itoa(i))

	// 	ioutil.ReadAll(resp.Body)

	// 	elapsed := time.Since(start)

	// 	times = append(times, elapsed.String())
	// }

	// var totalMs float64
	// for _, v := range times {
	// 	timeStringParts := strings.Split(v, "ms")
	// 	floatValue, _ := strconv.ParseFloat(timeStringParts[0], 64)
	// 	totalMs += floatValue
	// }

	// fmt.Println(totalMs / float64(len(times)))
	// fmt.Println(totalMs / float64(len(times)) * 7)
	// fmt.Println(totalMs / float64(len(times)) * 7 * 898 / 1000 / 60)

	// server.Start()
	go server.Start()
	//figure out how to use on vps lol (gives bad gateway with reverse proxy, maybe split in seperate projects)
	commands.Scanner()
}
