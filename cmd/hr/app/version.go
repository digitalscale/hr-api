package app

import (
	"runtime"

	"github.com/spf13/cobra"
)

func Version(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display server version.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("HR API Server %s %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
		},
	}
}
