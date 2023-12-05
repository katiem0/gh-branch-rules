package cmd

import (
	"github.com/spf13/cobra"

	createCmd "github.com/katiem0/gh-branch-rules/cmd/create"
	listCmd "github.com/katiem0/gh-branch-rules/cmd/list"
)

func NewCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "branchrules <command> [flags]",
		Short: "List and create branch protection rules.",
		Long:  "List and create branch protection rules.",
	}

	cmd.AddCommand(listCmd.NewCmdList())
	cmd.AddCommand(createCmd.NewCmdCreate())
	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	return cmd
}
