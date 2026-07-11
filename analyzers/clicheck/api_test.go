package clicheck_test

import (
	"context"
	"testing"

	"github.com/livebud/cli"
)

// TestAPIContract verifies that the methods we document in the template
// actually exist with the expected signatures. If this test fails, update
// the template in clicheck.go.
func TestAPIContract(t *testing.T) {
	c := cli.New("test", "test cli")

	// cli.New returns *CLI with Command, Flag, Arg, Args, Run, Parse
	cmd := c.Command("sub", "a subcommand")

	// Command returns something with the same interface
	_ = cmd.Command("nested", "nested subcommand")

	// Flag(name, help) returns *Flag with chainable methods
	var s string
	cmd.Flag("name", "help").String(&s)

	var b bool
	cmd.Flag("name", "help").Bool(&b).Default(false)

	cmd.Flag("name", "help").Enum(&s, "a", "b").Default("a")

	var opt *string
	cmd.Flag("name", "help").Optional().String(&opt)

	// Arg(name, help) returns *Arg with String
	cmd.Arg("name", "help").String(&s)

	// Args(name, help) returns *Args with Strings
	var ss []string
	cmd.Args("name", "help").Strings(&ss)

	// Run takes a func(ctx) error
	cmd.Run(func(ctx context.Context) error { return nil })

	// Parse takes context and args
	_ = c.Parse(context.Background())
}
