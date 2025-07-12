package main

import (
	"fmt"
	"github.com/epenick123/blogagg/internal/config" // <--- This line is key!
	"log"
	"os"
)

// Define your structs
type state struct {
	// your config pointer
	ptr *config.Config
}

type command struct {
	// name and args fields
	name string
	args []string
}

type commands struct {
	// your map of handlers
	cmd_names map[string]func(*state, command) error
}

// Define your methods and functions
func handlerLogin(s *state, cmd command) error {
	// implementation
	if len(cmd.args) == 0 {
		return fmt.Errorf("The login handler expects a single argument, the username.")
	}
	s.ptr.SetUser(cmd.args[0])
	fmt.Println("The user has been set.")

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	// implementation
	handler, exists := c.cmd_names[cmd.name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	// implementation
	c.cmd_names[name] = f
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	s2 := state{ptr: &cfg}

	cmds := commands{
		cmd_names: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("not enough arguments provided")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}

	if err := cmds.run(&s2, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
