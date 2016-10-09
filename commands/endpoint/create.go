package endpoint

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	cli "github.com/mariusmagureanu/cli_poc/commands"
)

type CreateEndpoint struct {
	cli.Command

	endpoint cli.SimpleEndpoint

	name *string
	host *string
	path *string
}

func NewCreateEndpoint() CreateEndpoint {
	var cc = CreateEndpoint{}

	cc.endpoint = cli.SimpleEndpoint{}
	cc.Arg1 = cli.CREATE_COMMAND
	cc.Arg2 = cli.ONE_ENDPOINT_ARG

	cc.Flagset = flag.NewFlagSet("endpoint", flag.ContinueOnError)
	cc.name = cc.Flagset.String("name", "", "Endpoint name. (Required)")
	cc.host = cc.Flagset.String("host", "", "Endpoint host. (Required)")
	cc.path = cc.Flagset.String("path", "", "Endpoint path. (Required)")
	return cc

}

func (c CreateEndpoint) GetFlagSet() *flag.FlagSet {
	return c.Flagset
}

func (c CreateEndpoint) GetArg1() string {
	return c.Arg1
}

func (c CreateEndpoint) GetArg2() string {
	return c.Arg2
}

func (c CreateEndpoint) Validate() error {
	var err error
	if err = c.Flagset.Parse(os.Args[3:]); err == nil {

		if *c.name == "" {
			return errors.New("Invalid endpoint name.")
		}

		if *c.host == "" {
			return errors.New("Invalid endpoint host.")
		}

		if *c.path == "" {
			return errors.New("Invalid endpoint path.")
		}
	}
	return err
}

func (c CreateEndpoint) Run() error {

	var err = c.Validate()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Usage:")
		c.Flagset.PrintDefaults()
		return err
	}

	c.Flagset.Visit(cli.FlagVisitor)
	status, err := cli.Do(http.MethodPost, "/endpoints", &c.endpoint, cli.BodyData)

	c.Output(status)
	return err
}

func (c CreateEndpoint) Output(status int) {

	switch status {
	case http.StatusCreated:
		fmt.Fprintln(cli.Writer, c.endpoint.ToString())
	case http.StatusBadRequest:
		fmt.Println("Invalid data for creating a new endpoint.")
	case http.StatusConflict:
		fmt.Printf("Endpoint %s already exists.\n", c.endpoint.Name)
	default:
		fmt.Println(status)
	}

	cli.Writer.Flush()
}
