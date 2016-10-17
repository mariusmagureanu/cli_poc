package main

import (
	"fmt"
	"github.com/mariusmagureanu/cli_poc/commands"
	"github.com/mariusmagureanu/cli_poc/commands/consumer"
	"github.com/mariusmagureanu/cli_poc/commands/endpoint"
	"github.com/mariusmagureanu/cli_poc/commands/module"
	"os"
	"strings"
)

var (
	err error

	showConsumers   = consumer.NewShowConsumers()
	showOneConsumer = consumer.NewShowOneConsumer()
	createConsumer  = consumer.NewCreateConsumer()
	deleteConsumer  = consumer.NewDeleteConsumer()

	showEndpoints   = endpoint.NewShowEndpoints()
	showOneEndpoint = endpoint.NewShowOneEndpoint()
	createEndpoint  = endpoint.NewCreateEndpoint()
	deleteEndpoint  = endpoint.NewDeleteEndpoint()

	showAllEndpointModules   = module.NewShowEndpointModuless()
	showOneEndpointModule    = module.NewShowEndpointModule()
	addModuleToEndpoint      = module.NewAddEndpointModule()
	removeModuleFromEndpoint = module.NewRemoveEndpointModule()

	cliCommands = []commands.Runner{createConsumer, showConsumers, deleteConsumer, showOneConsumer,
		showEndpoints, showOneEndpoint, createEndpoint, deleteEndpoint,
		showAllEndpointModules, showOneEndpointModule, addModuleToEndpoint, removeModuleFromEndpoint}
)

func showUsage() {
	fmt.Println("show | create | delete | update | add |  sub-command is required.")
	for _, command := range cliCommands {
		command.GetFlagSet().PrintDefaults()
	}
}

func main() {

	commands.SetHost("http://127.0.0.1")
	commands.SetPort("8089")

	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	var arg1 = os.Args[1]
	var arg2 = os.Args[2]

	for _, command := range cliCommands {
		if strings.EqualFold(arg1, command.GetArg1()) && strings.EqualFold(arg2, command.GetArg2()) {
			err = command.Run()
			break
		}
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
