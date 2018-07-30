package alias_config

import (
  "fmt"
  "path"

  genericConfig "github.com/tkporter/kubectl-common/src/generic-config"
  kubectlManager "github.com/tkporter/kubectl-common/src/kubectl-manager"

  "github.com/spf13/viper"
)

// Finds the version corresponding to an alias found in the alias config.
func GetConfigurationForAlias(aliasConfigPath, alias string) (version, kubeconfig string, err error) {
  aliasConfig, err := LoadConfig(aliasConfigPath)
  if err != nil {
    return "", "", err
  }
  aliases := aliasConfig.GetStringMap("aliases")
  if configurationInterface, ok := aliases[alias]; ok {
    configuration := configurationInterface.(map[string]interface{})
    return configuration["version"].(string), configuration["kubeconfig"].(string), nil
  } else {
    return "", "", fmt.Errorf("Alias %s not found in the alias config", alias)
  }
}

// Applies the alias and kubectl versions found in the alias config.
// Calls functions to ensure the proper kubectl versions are downloaded.
// Returns the directory of the config file and any errors
func ApplyAliasConfig(aliasConfigPath string) (string, error) {
  aliasConfig, err := LoadConfig(aliasConfigPath)
  if err != nil {
    return "", err
  }
  aliases := aliasConfig.GetStringMap("aliases")
  applyAliases(aliases)
  return path.Dir(aliasConfig.ConfigFileUsed()), nil
}

// Loads the alias config
func LoadConfig(aliasConfigPath string) (*viper.Viper, error) {
  aliasConfig := viper.New()
  // TODO change up from "." to something else
  if aliasConfigPath != "" {
    aliasConfig.AddConfigPath(aliasConfigPath)
  } else {
    dirPath, err := genericConfig.GetConfigDirPath()
    if err != nil {
      return nil, err
    }
    aliasConfig.AddConfigPath(dirPath)
  }
  aliasConfig.SetConfigName("alias-config")

  if err := aliasConfig.ReadInConfig(); err != nil {
    return nil, err
  }
  return aliasConfig, nil
}

// Ensures there are no repeats in the versions specified by the aliases,
// and sets up the kubectl versions.
func applyAliases(aliases map[string]interface{}) {
  versionMap := make(map[string]bool)
  for _, aliasConfigurationInterface := range aliases {
    aliasConfiguration := aliasConfigurationInterface.(map[string]interface{})
    versionMap[aliasConfiguration["version"].(string)] = true
  }
  kubectlManager.SetupKubectlVersions(versionMap)
}
