package cmd

import (
  "fmt"
  "os"

  internalConfig "github.com/tkporter/kubectl-common/src/internal-config"

  "github.com/spf13/cobra"
)

var cmdUseVersionAlias = &cobra.Command{
  Use:   "use-version-alias [alias]",
  Short: "Apply an alias config",
  Long: "Apply an alias config",
  Args: cobra.ExactArgs(1),
  Run: useVersionAlias,
}

func init() {
  RootCmd.AddCommand(cmdUseVersionAlias)
}

func useVersionAlias(cmd *cobra.Command, args []string) {
  alias := args[0]
  err := internalConfig.SetCurrentVersionAlias(alias)
  if err != nil {
    fmt.Println("Error:", err)
    os.Exit(1)
  }
  alias, version, err := internalConfig.GetCurrentVersionAlias()
  fmt.Printf("Now using alias %s for kubectl version %s\n", alias, version)
}
