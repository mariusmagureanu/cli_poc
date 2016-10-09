package endpoint

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	cli "github.com/mariusmagureanu/cli_poc/commands"
)

type ShowOneEndpoint struct {
	cli.Command
	name *string

	endpoint cli.SimpleEndpoint
}

func NewShowOneEndpoint() ShowOneEndpoint {
	var dc = ShowOneEndpoint{}

	dc.endpoint = cli.SimpleEndpoint{}

	dc.Flagset = flag.NewFlagSet("endpoint", flag.ContinueOnError)
	dc.name = dc.Flagset.String("name", "", "Endpoint name. (Required)")

	dc.Arg1 = cli.SHOW_COMMAND
	dc.Arg2 = cli.ONE_ENDPOINT_ARG
	return dc
}

func (c ShowOneEndpoint) GetFlagSet() *flag.FlagSet {
	return c.Flagset
}

func (c ShowOneEndpoint) GetArg1() string {
	return c.Arg1
}

func (c ShowOneEndpoint) GetArg2() string {
	return c.Arg2
}

func (c ShowOneEndpoint) Validate() error {
	if err := c.Flagset.Parse(os.Args[3:]); err == nil {
		if *c.name == "" {
			return errors.New("Invalid endpoint name.")
		}
	}
	return nil
}

func (c ShowOneEndpoint) Run() error {
	var err = c.Validate()
	var getUrl = fmt.Sprintf("/%s/%s", "endpoints", *c.name)
	status, err := cli.Do(http.MethodGet, getUrl, &c.endpoint, nil)

	c.Output(status)
	return err
}

func (c ShowOneEndpoint) Output(status int) {

	switch status {
	case http.StatusOK:
		fmt.Fprintln(cli.Writer, c.endpoint.ToString())
	case http.StatusNotFound:
		fmt.Printf("Endpoint %s not found.\n", *c.name)

	default:
		fmt.Println(status)
	}

	cli.Writer.Flush()
}
