package cli

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/onsattia/docs-drift-detector/internal/registry"
	"github.com/spf13/cobra"
)

// errReported marks a failure whose details were already printed, so Execute adds no usage text.
var errReported = errors.New("reported")

// Execute runs the docs-drift CLI and returns the process exit code.
func Execute(args []string, stdout, stderr io.Writer, version string) int {
	root := newRootCmd(version)
	root.SetArgs(args)
	root.SetOut(stdout)
	root.SetErr(stderr)

	err := root.Execute()
	switch {
	case err == nil:
		return 0
	case errors.Is(err, errReported):
		return 1
	default:
		// Cobra prints usage to stdout on error; usage belongs on stderr with the error.
		// If stderr itself fails there is nowhere left to report it, so the exit code is the signal.
		_, _ = fmt.Fprintln(stderr, err)
		_, _ = fmt.Fprint(stderr, root.UsageString())
		return 1
	}
}

func newRootCmd(version string) *cobra.Command {
	root := &cobra.Command{
		Use:           "docs-drift",
		Short:         "Detect drift between documentation and code",
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
		// NoArgs makes an unknown subcommand an error. Without it, cobra treats
		// a command with no subcommands as accepting any arguments.
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	root.AddCommand(newRegistryCmd())
	return root
}

func newRegistryCmd() *cobra.Command {
	registryCmd := &cobra.Command{
		Use:   "registry",
		Short: "Work with the sources.yaml registry",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	registryCmd.AddCommand(newRegistryValidateCmd())
	return registryCmd
}

func newRegistryValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <path>",
		Short: "Validate a sources.yaml registry and print a summary",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reg, err := registry.Load(args[0])
			if err != nil {
				printErrors(cmd.ErrOrStderr(), err)
				return errReported
			}
			printSummary(cmd.OutOrStdout(), reg)
			return nil
		},
	}
}

// printErrors writes each validation problem on its own line, so every error is visible in one run.
func printErrors(w io.Writer, err error) {
	var verr *registry.ValidationError
	if errors.As(err, &verr) {
		for _, e := range verr.Errors {
			_, _ = fmt.Fprintf(w, "error: %s\n", e)
		}
		return
	}
	_, _ = fmt.Fprintf(w, "error: %v\n", err)
}

func printSummary(w io.Writer, reg *registry.Registry) {
	_, _ = fmt.Fprintf(w, "Registry is valid: %d project(s)\n", len(reg.Projects))
	names := make([]string, 0, len(reg.Projects))
	for name := range reg.Projects {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		_, _ = fmt.Fprintf(w, "  %s\n", name)
		for _, source := range reg.Projects[name].Code {
			_, _ = fmt.Fprintf(w, "    - %s@%s\n", source.Repo, source.Ref)
		}
	}
}
