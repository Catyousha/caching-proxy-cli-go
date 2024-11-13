package cmd

import (
	"catyousha/caching-proxy/internal"

	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	var port int
	var origin string
	var clearAll bool
	cmd := &cobra.Command{
		Use: "caching-proxy --port <number> --origin <url>",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if clearAll {
				internal.ClearCache()
				return
			}
			internal.SetupProxy(port, origin)
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to listen on")
	cmd.Flags().StringVarP(&origin, "origin", "o", "", "Origin server URL") 
	cmd.Flags().BoolVar(&clearAll, "clear-all", false, "Clear all cached data")
	return cmd
}