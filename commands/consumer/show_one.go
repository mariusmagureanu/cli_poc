package consumer

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	cli "github.com/mariusmagureanu/cli_poc/commands"
	"strings"
)

type ShowOneConsumer struct {
	cli.Command
	name *string

	consumer cli.SimpleConsumer
}

func NewShowOneConsumer() ShowOneConsumer {
	var dc = ShowOneConsumer{}

	dc.consumer = cli.SimpleConsumer{}

	dc.Flagset = flag.NewFlagSet("show one consumer", flag.ExitOnError)
	dc.name = dc.Flagset.String("name", "", "Consumer name. (Required)")

	dc.Arg1 = cli.SHOW_COMMAND
	dc.Arg2 = cli.ONE_CONSUMER_ARG

	dc.Flagset.Usage = func() {
		fmt.Println("Show a specific consumer:")
		fmt.Println("show consumer [options]")
		dc.Flagset.PrintDefaults()
	}

	return dc
}

func (c ShowOneConsumer) GetFlagSet() *flag.FlagSet {
	return c.Flagset
}

func (c ShowOneConsumer) GetArg1() string {
	return c.Arg1
}

func (c ShowOneConsumer) GetArg2() string {
	return c.Arg2
}

func (c ShowOneConsumer) Validate() error {
	var err error
	if err = c.Flagset.Parse(os.Args[3:]); err == nil {
		if strings.TrimSpace(*c.name) == "" {
			return errors.New("Invalid consumer name.")
		}
	}
	return err
}

func (c ShowOneConsumer) Run() error {
	var err = c.Validate()

	if err != nil {
		return err
	}

	var getUrl = fmt.Sprintf("/%s/%s", "consumers", *c.name)
	status, err := cli.Do(http.MethodGet, getUrl, &c.consumer, nil)

	c.Output(status)
	return err
}

func (c ShowOneConsumer) Output(status int) {

	switch status {
	case http.StatusOK:
		fmt.Fprintln(cli.Writer, c.consumer.ToString())
	case http.StatusNotFound:
		fmt.Printf("Consumer %s not found.\n", *c.name)

	default:
		fmt.Println(status)
	}

	cli.Writer.Flush()
}
