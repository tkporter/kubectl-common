package cmd

import (
  "fmt"
  "os"

  internalConfig "github.com/tkporter/kubectl-common/src/internal-config"
  kubectlManager "github.com/tkporter/kubectl-common/src/kubectl-manager"

  "github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
  Use:   "kubectl-common",
  Short: "A wrapper of kubectl for managing different client versions.",
  Long: "A wrapper of kubectl for managing different client versions.",
  Args: cobra.ArbitraryArgs,
  Run: runKubectlCommand,
}

func Execute() {
  if err := RootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

// Runs a command with the previously configured version of kubectl.
func runKubectlCommand(cmd *cobra.Command, args []string) {
  // Get the alias & version that was configured earlier
  alias, version, err := internalConfig.GetCurrentVersionAlias()
  if err != nil {
    fmt.Println("Error getting alias version", err)
    os.Exit(1)
  }
  fmt.Printf("Using alias %s for kubectl version %s\n", alias, version)
  kubectlManager.RunKubectlCommand(version, args)
}
