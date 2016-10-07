package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/varnish/policy-engine/kvstore"
	"net/http"
	"os"
	"text/tabwriter"
)

const host = "http://localhost:8089/"

var (
	consumerFlagSet = flag.NewFlagSet("consumer", flag.ContinueOnError)
	endpointFlagSet = flag.NewFlagSet("endpoint", flag.ContinueOnError)

	consumerName  = consumerFlagSet.String("name", "", "Consumer name")
	consumerEMail = consumerFlagSet.String("email", "", "Consumer email")

	endpointHost = endpointFlagSet.String("host", "", "Endpoint host")
	endpointPath = endpointFlagSet.String("path", "", "Endpoint path")
	endpointName = endpointFlagSet.String("name", "", "Endpoint name")
)

func doRequest(method string, url string, out interface{}, body map[string]string) (int, error) {

	var (
		client      http.Client
		err         error
		bodyAsBytes []byte
		status      int
	)

	if body != nil {
		bodyAsBytes, err = json.Marshal(body)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	req, err := http.NewRequest(method, host+url, bytes.NewBuffer(bodyAsBytes))
	req.Close = true

	if err != nil {
		return http.StatusInternalServerError, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if out != nil {
		var decoder = json.NewDecoder(resp.Body)
		status = resp.StatusCode
		err = decoder.Decode(out)
		resp.Body.Close()
	}

	return status, err
}

func main() {

	var bodyData = make(map[string]string)
	var padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

	visitor := func(a *flag.Flag) {
		bodyData[a.Name] = a.Value.String()
	}

	switch os.Args[1] {
	case "create":

		switch os.Args[2] {
		case "consumer":
			if err := consumerFlagSet.Parse(os.Args[3:]); err == nil {
				consumerFlagSet.Visit(visitor)
				doRequest(http.MethodPost, "consumers", nil, bodyData)
			}
		case "endpoint":
			if err := consumerFlagSet.Parse(os.Args[3:]); err == nil {
				endpointFlagSet.Visit(visitor)
				doRequest(http.MethodPost, "endpoints", nil, bodyData)
			}
		}

	case "show":
		switch os.Args[2] {

		case "consumer":
			if err := consumerFlagSet.Parse(os.Args[3:]); err == nil {
				var consumer kvstore.Consumer
				status, _ := doRequest(http.MethodGet, "consumers/"+(*consumerName), &consumer, nil)
				fmt.Printf("%v", consumer)
				fmt.Println(status)
			}
		case "consumers":
			var consumers []kvstore.Consumer
			doRequest(http.MethodGet, "consumers", &consumers, nil)
			for _, c := range consumers {
				fmt.Fprintln(w, c.Name+"\t"+c.Email+"\t")
			}

		case "endpoint":
			if err := endpointFlagSet.Parse(os.Args[3:]); err == nil {
				var endpoint kvstore.Endpoint
				doRequest(http.MethodGet, "endpoints/"+(*endpointName), &endpoint, nil)
				fmt.Printf("%v", endpoint)
			}

		case "endpoints":
			var endpoints []kvstore.Endpoint
			doRequest(http.MethodGet, "endpoints", &endpoints, nil)
			for _, e := range endpoints {
				fmt.Fprintln(w, e.Host+"\t"+e.Path+"\t"+e.Name)
			}

		}
	}
	w.Flush()

}
