package interal_config

import (
  "bytes"
  "fmt"
  "io"
  "os"
  "path"

  aliasConfig "github.com/tkporter/kubectl-common/alias-config"
  genericConfig "github.com/tkporter/kubectl-common/generic-config"

  "github.com/spf13/viper"
)

const configFileName = "config"
var dirFileMode = os.FileMode(0755)
var defaultConfig = []byte(`{
  "aliasConfigPath": "",
  "current": {
    "alias": "",
    "version": "",
    "kubeconfig": ""
  }
}
`)

// Gets the path to the alias config found in the internal config.
func GetAliasConfigPath() (string, error) {
  config, err := LoadConfig()
  if err != nil {
    return "", err
  }
  return config.GetString("aliasConfigPath"), nil
}

// Sets the path to the alias config in the internal config.
func SetAliasConfigPath(aliasConfigPath string) error {
  config, err := LoadConfig()
  if err != nil {
    return err
  }
  config.Set("aliasConfigPath", aliasConfigPath)
  config.WriteConfig()
  return nil
}

// Gets the currently used alias & corresponding version
func GetCurrentConfiguration() (alias, version, kubeconfig string, err error) {
  config, err := LoadConfig()
  if err != nil {
    return "", "", "", err
  }
  current := config.GetStringMapString("current")
  return current["alias"], current["version"], current["kubeconfig"], nil
}

// Sets the new alias & version to use
func SetCurrentConfiguration(alias string) error {
  aliasConfigPath, err := GetAliasConfigPath()
  if err != nil {
    return err
  }
  version, kubeconfig, err := aliasConfig.GetConfigurationForAlias(aliasConfigPath, alias)
  if err != nil {
    return err
  }
  config, err := LoadConfig()
  if err != nil {
    return err
  }
  config.Set("current.alias", alias)
  config.Set("current.version", version)
  config.Set("current.kubeconfig", kubeconfig)
  config.WriteConfig()
  return nil
}

// Loads the internal config
func LoadConfig() (*viper.Viper, error) {
  dirPath, entirePath, err := getConfigPath()
  if err != nil {
    return nil, err
  }

  config := viper.New()
  config.SetConfigType("json")
  config.AddConfigPath(dirPath)
  config.SetConfigName(configFileName)
  if err := config.ReadInConfig(); err != nil {
    // If the config just isn't found, then create one
    _, ok := err.(viper.ConfigFileNotFoundError)
    if ok {
      err := createDefaultConfig(dirPath, entirePath)
      if err != nil {
        return nil, err
      }
      // try reading again...
      err = config.ReadInConfig()
      return config, err
    } else {
      return nil, err
    }
  }
  return config, nil
}

// Intended to be used if the config does not exist. Creates the necessary
// path to the config and initializes the config as the defaultConfig.
func createDefaultConfig(dirPath, entirePath string) error {
  // make sure the file path is there
  err := os.MkdirAll(dirPath, dirFileMode)
  if err != nil {
    return err
  }
  // then create the file
  out, err := os.Create(entirePath)
  if err != nil {
    return err
  }
  defer out.Close()

  // and copy the default content over
  defaultReader := bytes.NewReader(defaultConfig)
  _, err = io.Copy(out, defaultReader)
  return err
}

// Gets the path of the internal config
func getConfigPath() (dirPath string, entirePath string, err error) {
  dirPath, err = genericConfig.GetConfigDirPath()
  if err != nil {
    return "", "", err
  }
  entirePath = path.Join(dirPath, fmt.Sprintf("%s.json", configFileName))
  return dirPath, entirePath, nil
}
