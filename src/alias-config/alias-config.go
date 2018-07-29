package alias_config

import (
  "fmt"

  kubectlManager "github.com/tkporter/kubectl-common/src/kubectl-manager"

  "github.com/spf13/viper"
)

func GetVersionForAlias(aliasConfigPath, alias string) (string, error) {
  aliasConfig, err := LoadConfig(aliasConfigPath)
  if err != nil {
    return "", err
  }
  aliases := aliasConfig.GetStringMapString("aliases")
  if version, ok := aliases[alias]; ok {
    return version, nil
  } else {
    return "", fmt.Errorf("Alias %s not found in config", alias)
  }
}

func ApplyAliasConfig(aliasConfigPath string) (string, error) {
  aliasConfig, err := LoadConfig(aliasConfigPath)
  if err != nil {
    return "", err
  }
  aliases := aliasConfig.GetStringMapString("aliases")
  applyAliases(aliases)
  return aliasConfig.ConfigFileUsed(), nil
}

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

func applyAliases(aliases map[string]string) {
  versionMap := make(map[string]bool)
  for _, version := range aliases {
    versionMap[version] = true
  }
  kubectlManager.SetupKubectlVersions(versionMap)
}
