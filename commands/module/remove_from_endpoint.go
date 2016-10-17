package module

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	cli "github.com/mariusmagureanu/cli_poc/commands"
)

type RemoveEndpointModule struct {
	cli.Command
	endpointName *string
	moduleName   *string
}

func NewRemoveEndpointModule() RemoveEndpointModule {
	var sc = RemoveEndpointModule{}

	sc.Flagset = flag.NewFlagSet("modules", flag.ContinueOnError)
	sc.endpointName = sc.Flagset.String("endpoint", "", "Endpoint name. (Required)")
	sc.moduleName = sc.Flagset.String("name", "", "Module name. (Required)")

	sc.Arg1 = cli.DELETE_COMMAND
	sc.Arg2 = cli.ONE_MODULE_ARG

	return sc
}

func (c RemoveEndpointModule) GetFlagSet() *flag.FlagSet {
	return c.Flagset
}

func (c RemoveEndpointModule) GetArg1() string {
	return c.Arg1
}

func (c RemoveEndpointModule) GetArg2() string {
	return c.Arg2
}

func (c RemoveEndpointModule) Validate() error {
	var err error
	if err = c.Flagset.Parse(os.Args[3:]); err == nil {
		if strings.TrimSpace(*c.endpointName) == "" {
			return errors.New("Invalid endpoint name.")
		}
		if strings.TrimSpace(*c.moduleName) == "" {
			return errors.New("Invalid module name.")
		}
	}
	return err
}

func (c RemoveEndpointModule) Run() error {
	var err = c.Validate()

	if err != nil {
		return err
	}

	var removeModuleUrl = fmt.Sprintf("/%s/%s/%s/%s", "endpoints", *c.endpointName, "modules", *c.moduleName)
	status, err := cli.Do(http.MethodDelete, removeModuleUrl, nil, nil)

	c.Output(status)
	return err
}

func (c RemoveEndpointModule) Output(status int) {

	switch status {
	case http.StatusNoContent:
		fmt.Printf("Module %s has been removed fron endpoint %s.\n", *c.moduleName, *c.endpointName)

	default:
		fmt.Println(status)
	}

	cli.Writer.Flush()
}
