package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	response, err := http.DefaultClient.Get("http://127.0.0.1:8000/a?echo=hello")
	if err != nil {
		log.Fatalln(err)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("client go data:%s", data)
}
