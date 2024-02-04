package main

import (
	"fmt"
	"log"
	"net/http"
)

func handlerErr(err error, chanErr chan error) {
	if err != nil {
		chanErr <- err
	}
}

func httpGet(client *http.Client, api string, errChan chan error, resChan chan string, service string) {
	_, err := client.Get(api)
	handlerErr(err, errChan)
	resChan <- service
}

func main() {
	errChan := make(chan error)
	resChan := make(chan string)
	client := http.DefaultClient
	cep := "01311200"

	brasilApi := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	viaApi := fmt.Sprintf("http://viacep.com.br/ws/%s/json", cep)

	go func() {
		for {
			httpGet(client, viaApi, errChan, resChan, "viaApi")
		}
	}()
	go func() {
		for {
			httpGet(client, brasilApi, errChan, resChan, "brasilApi")
		}
	}()

	for {
		select {
		case msg := <-resChan:
			fmt.Printf("Chegou primeiro o: %s\n", msg)
		case err := <-errChan:
			log.Fatalf("%+v", err)
		}
	}
}
