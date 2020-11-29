package app

import "github.com/spf13/cobra"

func NewDefaultCommand(version string) *cobra.Command {
	root := &cobra.Command{
		Use:   "hr",
		Short: "HR API server.",
	}

	root.AddCommand(Server())
	root.AddCommand(Version(version))

	return root
}
