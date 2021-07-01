package functions

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//
func GetRequest(url string) []byte {
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
