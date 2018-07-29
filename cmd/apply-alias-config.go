package cmd

import (
  "fmt"

  aliasConfig "github.com/tkporter/kubectl-common/src/alias-config"
  internalConfig "github.com/tkporter/kubectl-common/src/internal-config"

  "github.com/spf13/cobra"
)

var AliasConfigPath string

var cmdApplyAliasConfig = &cobra.Command{
  Use:   "apply-alias-config",
  Short: "Apply an alias config",
  Long: "Apply an alias config",
  Run: applyAliasConfig,
}

func init() {
  cmdApplyAliasConfig.Flags().StringVarP(&AliasConfigPath, "config-dir", "d", "", "The directory to read the alias config from")
  RootCmd.AddCommand(cmdApplyAliasConfig)
}

func applyAliasConfig(cmd *cobra.Command, args []string) {
  aliasConfigPath, err := aliasConfig.ApplyAliasConfig(AliasConfigPath)
  if err != nil {
    fmt.Println("Error applying alias config:", err)
  }
  err = internalConfig.SetAliasConfigPath(aliasConfigPath)
  if err != nil {
    fmt.Println("Error applying alias config:", err)
  }
}
