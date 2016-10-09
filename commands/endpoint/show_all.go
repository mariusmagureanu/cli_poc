package endpoint

import (
	"flag"
	"fmt"
	"net/http"

	cli "github.com/mariusmagureanu/cli_poc/commands"
)

type ShowEndpoints struct {
	cli.Command

	endpoints []cli.SimpleEndpoint
}

func NewShowEndpoints() ShowEndpoints {
	var sc = ShowEndpoints{}

	sc.Flagset = flag.NewFlagSet("endpoints", flag.ContinueOnError)
	sc.Arg1 = cli.SHOW_COMMAND
	sc.Arg2 = cli.ALL_ENDPOINTS_ARG

	return sc
}

func (c ShowEndpoints) GetFlagSet() *flag.FlagSet {
	return c.Flagset
}

func (c ShowEndpoints) GetArg1() string {
	return c.Arg1
}

func (c ShowEndpoints) GetArg2() string {
	return c.Arg2
}

func (c ShowEndpoints) Validate() error {
	return nil
}

func (c ShowEndpoints) Run() error {
	var err = c.Validate()
	status, err := cli.Do(http.MethodGet, "/endpoints", &c.endpoints, nil)

	c.Output(status)
	return err
}

func (c ShowEndpoints) Output(status int) {

	switch status {
	case http.StatusOK:
		for _, endpoint := range c.endpoints {
			fmt.Fprintln(cli.Writer, endpoint.ToString())
		}

	default:
		fmt.Println(status)
	}

	cli.Writer.Flush()
}
