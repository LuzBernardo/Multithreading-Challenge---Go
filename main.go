package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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

func handler(msg response, chanErr chan error) {
	fmt.Printf("Quem chegou primeiro foi o: %s\n", msg.api)

	reader, err := io.ReadAll(msg.res.Body)
	handlerErr(err, chanErr)
	defer msg.res.Body.Close()

	fmt.Printf("Com o seguinte body de response: %s\n", string(reader))
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

	select {
	case msg := <-resChan:
		handler(msg, errChan)
	case err := <-errChan:
		log.Fatalf("%+v", err)
	case <-time.After(time.Second):
		log.Fatalln("time out!!!")
	}
}
