package endpoint

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	cli "github.com/mariusmagureanu/cli_poc/commands"
	"strings"
)

type DeleteEndpoint struct {
	cli.Command
	name *string
}

func NewDeleteEndpoint() DeleteEndpoint {
	var dc = DeleteEndpoint{}

	dc.Flagset = flag.NewFlagSet("delete endpoint", flag.ExitOnError)
	dc.name = dc.Flagset.String("name", "", "Endpoint name. (Required)")

	dc.Arg1 = cli.DELETE_COMMAND
	dc.Arg2 = cli.ONE_ENDPOINT_ARG

	dc.Flagset.Usage = func() {
		fmt.Println("delete endpoint [options]")
		dc.Flagset.PrintDefaults()
	}

	return dc
}

func (c DeleteEndpoint) GetFlagSet() *flag.FlagSet {
	return c.Flagset
}

func (c DeleteEndpoint) GetArg1() string {
	return c.Arg1
}

func (c DeleteEndpoint) GetArg2() string {
	return c.Arg2
}

func (c DeleteEndpoint) Validate() error {
	var err error
	if err = c.Flagset.Parse(os.Args[3:]); err == nil {
		if strings.TrimSpace(*c.name) == "" {
			return errors.New("Invalid endpoint name.")
		}
	}
	return err
}

func (c DeleteEndpoint) Run() error {
	var err = c.Validate()

	if err != nil {
		return err
	}

	var deleteUrl = fmt.Sprintf("/%s/%s", "endpoints", *c.name)
	status, err := cli.Do(http.MethodDelete, deleteUrl, nil, nil)

	c.Output(status)
	return err
}

func (c DeleteEndpoint) Output(status int) {

	switch status {
	case http.StatusNoContent:
		fmt.Printf("Endpoint %s has been deleted.\n", *c.name)
	case http.StatusNotFound:
		fmt.Printf("Endpoint %s not found.\n", *c.name)

	default:
		fmt.Println(status)
	}

	cli.Writer.Flush()
}
