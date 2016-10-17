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

type AddEndpointModule struct {
	cli.Command
	endpointName *string
	moduleName   *string

	module map[string]interface{}
}

func NewAddEndpointModule() AddEndpointModule {
	var sc = AddEndpointModule{}

	sc.Flagset = flag.NewFlagSet("modules", flag.ExitOnError)
	sc.endpointName = sc.Flagset.String("endpoint", "", "Endpoint name. (Required)")
	sc.moduleName = sc.Flagset.String("name", "", "Module name. (Required)")

	sc.Arg1 = cli.ADD_COMMAND
	sc.Arg2 = cli.ONE_MODULE_ARG

	sc.Flagset.Usage = func() {
		fmt.Println("add module [name] endpoint [endpoint]")
		sc.Flagset.PrintDefaults()
	}

	return sc
}

func (c AddEndpointModule) GetFlagSet() *flag.FlagSet {
	return c.Flagset
}

func (c AddEndpointModule) GetArg1() string {
	return c.Arg1
}

func (c AddEndpointModule) GetArg2() string {
	return c.Arg2
}

func (c AddEndpointModule) Validate() error {
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

func (c AddEndpointModule) Run() error {
	var err = c.Validate()

	if err != nil {
		return err
	}

	c.Flagset.Visit(cli.FlagVisitor)
	var addModuleUrl = fmt.Sprintf("/%s/%s/%s", "endpoints", *c.endpointName, "modules")
	status, err := cli.Do(http.MethodPost, addModuleUrl, &c.module, cli.BodyData)

	c.Output(status)
	return err
}

func (c AddEndpointModule) Output(status int) {

	switch status {
	case http.StatusCreated:
		out, _ := json.MarshalIndent(c.module, "", "  ")
		fmt.Println(string(out))

	default:
		fmt.Println(status)
	}

	cli.Writer.Flush()
}
