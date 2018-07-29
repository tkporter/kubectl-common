package interal_config

import (
  "bytes"
  "fmt"
  "io"
  "os"
  "path"

  aliasConfig "github.com/tkporter/kubectl-common/src/alias-config"

  homedir "github.com/mitchellh/go-homedir"
  "github.com/spf13/viper"
)

const configFileName = "config"
var dirFileMode = os.FileMode(0755)
var defaultConfig = []byte(`{
  "aliasConfigPath": "",
  "current": {
    "alias": "",
    "version": ""
  }
}
`)

func GetAliasConfigPath() (string, error) {
  config, err := LoadConfig()
  if err != nil {
    return "", err
  }
  return config.GetString("current"), nil
}

func SetAliasConfigPath(aliasConfigPath string) error {
  config, err := LoadConfig()
  if err != nil {
    return err
  }
  config.Set("aliasConfigPath", aliasConfigPath)
  config.WriteConfig()
  return nil
}

func GetCurrentVersionAlias() (alias, version string, err error) {
  config, err := LoadConfig()
  if err != nil {
    return "", "", err
  }
  current := config.GetStringMapString("current")
  return current["alias"], current["version"], nil
}

func SetCurrentVersionAlias(alias string) error {
  aliasConfigPath, err := GetAliasConfigPath()
  if err != nil {
    return err
  }
  version, err := aliasConfig.GetVersionForAlias(aliasConfigPath, alias)
  if err != nil {
    return err
  }
  config, err := LoadConfig()
  if err != nil {
    return err
  }
  config.Set("current.alias", alias)
  config.Set("current.version", version)
  config.WriteConfig()
  return nil
}

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

func getConfigPath() (dirPath string, entirePath string, err error) {
  home, err := homedir.Dir()
  if err != nil {
    return "", "", err
  }
  dirPath = path.Join(home, ".kube", "kube-common")
  entirePath = path.Join(dirPath, fmt.Sprintf("%s.json", configFileName))
  return dirPath, entirePath, nil
}
