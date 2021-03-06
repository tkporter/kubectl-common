package cmd

import (
  "fmt"
  "os"

  internalConfig "github.com/tkporter/kubectl-common/pkg/internal-config"
  kubectlManager "github.com/tkporter/kubectl-common/pkg/kubectl-manager"

  "github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
  Use:   "kubectl-common",
  Short: "A wrapper of kubectl for managing different client versions.",
  Long: "A wrapper of kubectl for managing different client versions.",
  Args: cobra.ArbitraryArgs,
  Run: runKubectlCommand,
  DisableFlagParsing: true,
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
  _, version, kubeconfig, err := internalConfig.GetCurrentConfiguration()
  if err != nil {
    fmt.Println("Error getting alias configuration:", err)
    os.Exit(1)
  }
  err = kubectlManager.RunKubectlCommand(version, kubeconfig, args)
  if err != nil {
    fmt.Println("Error executing kubectl command:", err)
    os.Exit(1)
  }
}
