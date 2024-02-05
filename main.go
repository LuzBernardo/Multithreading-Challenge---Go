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

type request struct {
	client  *http.Client
	api     string
	errChan *chan error
	resChan *chan response
	service string
}

func handlerErr(err error, chanErr *chan error) {
	if err != nil {
		*chanErr <- err
	}
}

func handler(msg response, chanErr *chan error) {
	fmt.Printf("Quem chegou primeiro foi o: %s\n", msg.api)

	reader, err := io.ReadAll(msg.res.Body)
	handlerErr(err, chanErr)
	defer msg.res.Body.Close()

	fmt.Printf("Com o seguinte body de response: %s\n", string(reader))
}

func httpGet(req *request) {
	res, err := req.client.Get(req.api)
	handlerErr(err, req.errChan)
	response := response{
		res: res,
		api: req.service,
	}
	*req.resChan <- response
}

func main() {
	errChan := make(chan error)
	resChan := make(chan response, 2)
	client := http.DefaultClient
	cep := "01311200"

	brasilApi := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	viaApi := fmt.Sprintf("http://viacep.com.br/ws/%s/json", cep)

	go func() {
		request := request{
			client:  client,
			api:     viaApi,
			errChan: &errChan,
			resChan: &resChan,
			service: "viaApi",
		}
		httpGet(&request)
	}()
	go func() {
		request := request{
			client:  client,
			api:     brasilApi,
			errChan: &errChan,
			resChan: &resChan,
			service: "brasilApi",
		}
		httpGet(&request)
	}()

	select {
	case msg := <-resChan:
		handler(msg, &errChan)
	case err := <-errChan:
		log.Fatalf("%+v", err)
	case <-time.After(time.Second):
		log.Fatalln("time out!!!")
	}
}
