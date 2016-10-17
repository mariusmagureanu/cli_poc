package module

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	cli "github.com/mariusmagureanu/cli_poc/commands"
)

type ShowEndpointModule struct {
	cli.Command
	endpointName *string
	moduleName   *string

	module map[string]interface{}
}

func NewShowEndpointModule() ShowEndpointModule {
	var sc = ShowEndpointModule{}

	sc.Flagset = flag.NewFlagSet("modules", flag.ExitOnError)
	sc.endpointName = sc.Flagset.String("endpoint", "", "Endpoint name. (Required)")
	sc.moduleName = sc.Flagset.String("name", "", "Module name. (Required)")

	sc.Arg1 = cli.SHOW_COMMAND
	sc.Arg2 = cli.ONE_MODULE_ARG

	sc.Flagset.Usage = func() {
		fmt.Println("show module [name] endpoint [endpoint]")
		sc.Flagset.PrintDefaults()
	}
	return sc
}

func (c ShowEndpointModule) GetFlagSet() *flag.FlagSet {
	return c.Flagset
}

func (c ShowEndpointModule) GetArg1() string {
	return c.Arg1
}

func (c ShowEndpointModule) GetArg2() string {
	return c.Arg2
}

func (c ShowEndpointModule) Validate() error {
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

func (c ShowEndpointModule) Run() error {
	var err = c.Validate()

	if err != nil {
		return err
	}

	var showModuleUrl = fmt.Sprintf("/%s/%s/%s/%s", "endpoints", *c.endpointName, "modules", *c.moduleName)
	status, err := cli.Do(http.MethodGet, showModuleUrl, &c.module, nil)

	c.Output(status)
	return err
}

func (c ShowEndpointModule) Output(status int) {

	switch status {
	case http.StatusOK:
		out, _ := json.MarshalIndent(c.module, "", "  ")
		fmt.Println(string(out))

	default:
		fmt.Println(status)
	}

	cli.Writer.Flush()
}
