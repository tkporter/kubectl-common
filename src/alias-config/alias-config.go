package alias_config

import (
  "fmt"

  kubectlManager "github.com/tkporter/kubectl-common/src/kubectl-manager"

  "github.com/spf13/viper"
)

// Finds the version corresponding to an alias found in the alias config.
func GetVersionForAlias(aliasConfigPath, alias string) (string, error) {
  aliasConfig, err := LoadConfig(aliasConfigPath)
  if err != nil {
    return "", err
  }
  aliases := aliasConfig.GetStringMapString("aliases")
  if version, ok := aliases[alias]; ok {
    return version, nil
  } else {
    return "", fmt.Errorf("Alias %s not found in the alias config", alias)
  }
}

// Applies the alias and kubectl versions found in the alias config.
// Calls functions to ensure the proper kubectl versions are downloaded.
func ApplyAliasConfig(aliasConfigPath string) (string, error) {
  aliasConfig, err := LoadConfig(aliasConfigPath)
  if err != nil {
    return "", err
  }
  aliases := aliasConfig.GetStringMapString("aliases")
  applyAliases(aliases)
  return aliasConfig.ConfigFileUsed(), nil
}

// Loads the alias config
func LoadConfig(aliasConfigPath string) (*viper.Viper, error) {
  aliasConfig := viper.New()
  // TODO change up from "." to something else
  if aliasConfigPath != "" {
    aliasConfig.AddConfigPath(aliasConfigPath)
  } else {
    aliasConfig.AddConfigPath(".")
  }
  aliasConfig.SetConfigName("alias-config")

  if err := aliasConfig.ReadInConfig(); err != nil {
    return nil, err
  }
  return aliasConfig, nil
}

// Ensures there are no repeats in the versions specified by the aliases,
// and sets up the kubectl versions.
func applyAliases(aliases map[string]string) {
  versionMap := make(map[string]bool)
  for _, version := range aliases {
    versionMap[version] = true
  }
  kubectlManager.SetupKubectlVersions(versionMap)
}
