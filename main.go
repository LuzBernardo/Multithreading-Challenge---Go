package main

import (
	"fmt"
	"log"
	"net/http"
)

type response struct {
	res *http.Response
	api string
}

func handlerErr(err error, chanErr chan error) {
	if err != nil {
		chanErr <- err
	}
}

func httpGet(client *http.Client, api string, errChan chan error, resChan chan response, service string) {
	res, err := client.Get(api)
	handlerErr(err, errChan)
	response := response{
		res: res,
		api: service,
	}
	resChan <- response
}

func main() {
	errChan := make(chan error)
	resChan := make(chan response)
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
			fmt.Printf("Chegou primeiro o: %s\n", msg.api)
		case err := <-errChan:
			log.Fatalf("%+v", err)
		}
	}
}
