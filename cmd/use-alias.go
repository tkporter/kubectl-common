package cmd

import (
  "fmt"
  "os"

  internalConfig "github.com/tkporter/kubectl-common/internal-config"

  "github.com/spf13/cobra"
)

var cmdUseVersionAlias = &cobra.Command{
  Use:   "use-alias [alias]",
  Short: "Change to the configuration for an alias",
  Long: "Sets the current configuration to the alias's version and kubeconfig.",
  Args: cobra.ExactArgs(1),
  Run: useVersionAlias,
}

func init() {
  RootCmd.AddCommand(cmdUseVersionAlias)
}

// Configures kubectl-common to use the kubectl version corresponding to an
// alias found in the alias config.
func useVersionAlias(cmd *cobra.Command, args []string) {
  alias := args[0]
  err := internalConfig.SetCurrentConfiguration(alias)
  if err != nil {
    fmt.Println("Error setting version alias:", err)
    os.Exit(1)
  }
  alias, version, kubeconfig, err := internalConfig.GetCurrentConfiguration()
  fmt.Printf("Now using alias %s for kubectl version %s and kubeconfig %s\n", alias, version, kubeconfig)
}
