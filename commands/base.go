package commands

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"
)

const (
	CREATE_COMMAND = "create"
	SHOW_COMMAND   = "show"
	DELETE_COMMAND = "delete"
	UPDATE_COMMAND = "update"
	ADD_COMMAND    = "add"

	ONE_CONSUMER_ARG  = "consumer"
	ALL_CONSUMERS_ARG = "consumers"

	ONE_ENDPOINT_ARG  = "endpoint"
	ALL_ENDPOINTS_ARG = "endpoints"

	ONE_MODULE_ARG = "module"
	ALL_MODULE_ARG = "modules"

	padding = 3
)

var (
	host     string
	port     string
	BodyData map[string]string

	Writer *tabwriter.Writer

	FlagVisitor = func(a *flag.Flag) {
		BodyData[a.Name] = a.Value.String()
	}
)

func init() {
	// defaults, if not set otherwise
	host = "http://127.0.0.1"
	port = "8089"

	BodyData = make(map[string]string)
	Writer = tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
}

type SimpleConsumer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SimpleEndpoint struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Path string `json:"path"`
}

func (s *SimpleConsumer) ToString() string {
	return s.Name + "\t" + s.Email + "\t"
}

func (s *SimpleEndpoint) ToString() string {
	return s.Name + "\t" + s.Host + "\t" + s.Path + "\t"
}

func SetHost(h string) {
	host = h
}

func SetPort(p string) {
	port = p
}

type Command struct {
	Arg1    string
	Arg2    string
	Flagset *flag.FlagSet
}

type Runner interface {
	Validate() error
	GetArg1() string
	GetArg2() string
	GetFlagSet() *flag.FlagSet
	Run() error
	Output(int)
}

// Do performs a http request against the admin server.
//
// Input
// - Method: http method
// - Url: request url
// - Out: generic interface in which the response of the
//        request will be marshaled into.
// - Body: Key/value struct consisting out of cli flags,
//         these will be picked up and sent as request payload.
//
//Output
// - An integer which defines the status code of the request.
// - An error that may appear along the way.
func Do(method string, url string, out interface{}, body map[string]string) (int, error) {

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

	reqUrl := fmt.Sprintf("%s:%s%s", host, port, url)
	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(bodyAsBytes))

	if err != nil {
		return http.StatusInternalServerError, err
	}
	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	status = resp.StatusCode

	if out != nil {
		var decoder = json.NewDecoder(resp.Body)
		err = decoder.Decode(out)
		resp.Body.Close()
	}

	return status, err
}
