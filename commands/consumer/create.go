package consumer

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	cli "github.com/mariusmagureanu/cli_poc/commands"
)

type CreateConsumer struct {
	cli.Command

	consumer cli.SimpleConsumer

	name  *string
	email *string
}

func NewCreateConsumer() CreateConsumer {
	var cc = CreateConsumer{}

	cc.consumer = cli.SimpleConsumer{}
	cc.Arg1 = cli.CREATE_COMMAND
	cc.Arg2 = cli.ONE_CONSUMER_ARG

	cc.Flagset = flag.NewFlagSet("Create consumer", flag.ExitOnError)
	cc.name = cc.Flagset.String("name", "", "Consumer name. (Required)")
	cc.email = cc.Flagset.String("email", "", "Consumer email.")

	cc.Flagset.Usage = func() {
		fmt.Println("Createa a new consumer:")
		fmt.Println("create consumer [options]")
		cc.Flagset.PrintDefaults()
	}
	return cc

}

func (c CreateConsumer) GetFlagSet() *flag.FlagSet {
	return c.Flagset
}

func (c CreateConsumer) GetArg1() string {
	return c.Arg1
}

func (c CreateConsumer) GetArg2() string {
	return c.Arg2
}

func (c CreateConsumer) Validate() error {
	var err error

	if err = c.Flagset.Parse(os.Args[3:]); err == nil {
		if *c.name == "" {
			return errors.New("Invalid consumer name.")
		}

	}
	return err
}

func (c CreateConsumer) Run() error {

	var err = c.Validate()

	if err != nil {
		return err
	}

	c.Flagset.Visit(cli.FlagVisitor)
	status, err := cli.Do(http.MethodPost, "/consumers", &c.consumer, cli.BodyData)

	c.Output(status)
	return err
}

func (c CreateConsumer) Output(status int) {

	switch status {
	case http.StatusCreated:
		fmt.Fprintln(cli.Writer, c.consumer.ToString())
	case http.StatusBadRequest:
		fmt.Println("Invalid data for creating a new consumer.")
	case http.StatusConflict:
		fmt.Printf("Consumer %s already exists.\n", c.consumer.Name)
	default:
		fmt.Println(status)
	}

	cli.Writer.Flush()
}
