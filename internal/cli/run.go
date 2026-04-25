package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

// RegisterRunCmd registers the run command which executes a subprocess with
// the resolved environment from the given profile chain.
func RegisterRunCmd(app *cli.App, cmdr *Commander) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:      "run",
		Usage:     "Run a command with resolved env vars from a profile chain",
		ArgsUsage: "<profile> [profile...] -- <command> [args...]",
		Action: func(c *cli.Context) error {
			args := c.Args().Slice()
			sep := -1
			for i, a := range args {
				if a == "--" {
					sep = i
					break
				}
			}
			if sep < 0 || sep == len(args)-1 {
				return fmt.Errorf("usage: envchain run <profile...> -- <command> [args...]")
			}
			profiles := args[:sep]
			cmdArgs := args[sep+1:]
			return cmdr.RunExec(profiles, cmdArgs)
		},
	})
}

// RunExec resolves the profile chain and executes the given command with the
// merged environment variables injected into the process environment.
func (cmdr *Commander) RunExec(profiles []string, cmdArgs []string) error {
	if len(profiles) == 0 {
		return fmt.Errorf("at least one profile is required")
	}
	if len(cmdArgs) == 0 {
		return fmt.Errorf("no command specified")
	}

	merged := map[string]string{}
	for _, name := range profiles {
		p, err := cmdr.store.GetProfile(name)
		if err != nil {
			return fmt.Errorf("profile %q not found", name)
		}
		for k, v := range p.Vars {
			if _, exists := merged[k]; !exists {
				merged[k] = v
			}
		}
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = os.Environ()
	for k, v := range merged {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
