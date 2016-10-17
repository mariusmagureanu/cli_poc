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

type DeleteConsumer struct {
	cli.Command
	name *string
}

func NewDeleteConsumer() DeleteConsumer {
	var dc = DeleteConsumer{}

	dc.Flagset = flag.NewFlagSet("delete consumer", flag.ExitOnError)
	dc.name = dc.Flagset.String("name", "", "Consumer name. (Required)")

	dc.Arg1 = cli.DELETE_COMMAND
	dc.Arg2 = cli.ONE_CONSUMER_ARG

	dc.Flagset.Usage = func() {
		fmt.Println("delete consumer [options]")
		dc.Flagset.PrintDefaults()
	}

	return dc
}

func (c DeleteConsumer) GetFlagSet() *flag.FlagSet {
	return c.Flagset
}

func (c DeleteConsumer) GetArg1() string {
	return c.Arg1
}

func (c DeleteConsumer) GetArg2() string {
	return c.Arg2
}

func (c DeleteConsumer) Validate() error {
	var err error
	if err = c.Flagset.Parse(os.Args[3:]); err == nil {
		if strings.TrimSpace(*c.name) == "" {
			return errors.New("Invalid consumer name.")
		}
	}
	return err
}

func (c DeleteConsumer) Run() error {
	var err = c.Validate()

	if err != nil {
		return err
	}

	var deleteUrl = fmt.Sprintf("/%s/%s", "consumers", *c.name)
	status, err := cli.Do(http.MethodDelete, deleteUrl, nil, nil)

	c.Output(status)
	return err
}

func (c DeleteConsumer) Output(status int) {

	switch status {
	case http.StatusNoContent:
		fmt.Printf("Consumer %s has been deleted.\n", *c.name)
	case http.StatusNotFound:
		fmt.Printf("Consumer %s not found.\n", *c.name)

	default:
		fmt.Println(status)
	}

	cli.Writer.Flush()
}
