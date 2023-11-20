package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	RootCmd = &cobra.Command{
		Use:   "ysl auto",
		Short: "ysl",
		Long:  "ysl web tool",
	}
)

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8081, "Port to listen on")
}
