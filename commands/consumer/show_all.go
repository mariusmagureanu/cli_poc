package consumer

import (
	"flag"
	"fmt"
	"net/http"

	cli "github.com/mariusmagureanu/cli_poc/commands"
)

type ShowConsumers struct {
	cli.Command

	consumers []cli.SimpleConsumer
}

func NewShowConsumers() ShowConsumers {
	var sc = ShowConsumers{}

	sc.Flagset = flag.NewFlagSet("show consumers", flag.ExitOnError)
	sc.Arg1 = cli.SHOW_COMMAND
	sc.Arg2 = cli.ALL_CONSUMERS_ARG

	sc.Flagset.Usage = func() {
		fmt.Println("Show all consumers:")
		fmt.Println("show consumers")
		sc.Flagset.PrintDefaults()
	}

	return sc
}

func (c ShowConsumers) GetFlagSet() *flag.FlagSet {
	return c.Flagset
}

func (c ShowConsumers) GetArg1() string {
	return c.Arg1
}

func (c ShowConsumers) GetArg2() string {
	return c.Arg2
}

func (c ShowConsumers) Validate() error {
	return nil
}

func (c ShowConsumers) Run() error {
	var err = c.Validate()
	status, err := cli.Do(http.MethodGet, "/consumers", &c.consumers, nil)

	c.Output(status)
	return err
}

func (c ShowConsumers) Output(status int) {

	switch status {
	case http.StatusOK:
		for _, consumer := range c.consumers {
			fmt.Fprintln(cli.Writer, consumer.ToString())
		}

	default:
		fmt.Println(status)
	}

	cli.Writer.Flush()
}
