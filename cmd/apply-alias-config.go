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
  Long: "Applies an alias config and downloads any necessary versions of kubectl.",
  Run: applyAliasConfig,
}

func init() {
  cmdApplyAliasConfig.Flags().StringVarP(&AliasConfigPath, "config-dir", "d", "", "The directory to read the alias config from")
  RootCmd.AddCommand(cmdApplyAliasConfig)
}

// Gets all kubectl versions that are needed and saves the path of the
// alias config to the internal config.
func applyAliasConfig(cmd *cobra.Command, args []string) {
  // Get all kubectl versions we need
  aliasConfigPath, err := aliasConfig.ApplyAliasConfig(AliasConfigPath)
  if err != nil {
    fmt.Println("Error applying alias config:", err)
  }
  // Save the path of the alias config to the internal config
  err = internalConfig.SetAliasConfigPath(aliasConfigPath)
  if err != nil {
    fmt.Println("Error saving the alias config path:", err)
  }
}
