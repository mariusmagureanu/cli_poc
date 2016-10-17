package module

import (
	"flag"
	"fmt"
	"net/http"

	"encoding/json"
	"errors"
	cli "github.com/mariusmagureanu/cli_poc/commands"
	"os"
	"strings"
)

type ShowEndpointModules struct {
	cli.Command
	endpointName *string

	modules []interface{}
}

func NewShowEndpointModuless() ShowEndpointModules {
	var sc = ShowEndpointModules{}

	sc.Flagset = flag.NewFlagSet("modules", flag.ExitOnError)
	sc.endpointName = sc.Flagset.String("endpoint", "", "Endpoint name. (Required)")

	sc.Arg1 = cli.SHOW_COMMAND
	sc.Arg2 = cli.ALL_MODULE_ARG

	sc.Flagset.Usage = func() {
		fmt.Println("show modules [endpoint]")
		sc.Flagset.PrintDefaults()
	}
	return sc
}

func (c ShowEndpointModules) GetFlagSet() *flag.FlagSet {
	return c.Flagset
}

func (c ShowEndpointModules) GetArg1() string {
	return c.Arg1
}

func (c ShowEndpointModules) GetArg2() string {
	return c.Arg2
}

func (c ShowEndpointModules) Validate() error {
	var err error
	if err = c.Flagset.Parse(os.Args[3:]); err == nil {
		if strings.TrimSpace(*c.endpointName) == "" {
			return errors.New("Invalid endpoint name.")
		}
	}
	return err
}

func (c ShowEndpointModules) Run() error {
	var err = c.Validate()
	if err != nil {
		return err
	}

	var showModulesUrl = fmt.Sprintf("/%s/%s/%s", "endpoints", *c.endpointName, "modules")
	status, err := cli.Do(http.MethodGet, showModulesUrl, &c.modules, nil)

	c.Output(status)
	return err
}

func (c ShowEndpointModules) Output(status int) {

	switch status {
	case http.StatusOK:
		out, _ := json.MarshalIndent(c.modules, "", "  ")
		fmt.Println(string(out))

	default:
		fmt.Println(status)
	}

	cli.Writer.Flush()
}
