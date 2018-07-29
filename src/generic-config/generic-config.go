package generic_config

import (
  "path"
  
  homedir "github.com/mitchellh/go-homedir"
)

func GetConfigDirPath() (string, error) {
  home, err := homedir.Dir()
  if err != nil {
    return "", err
  }
  dirPath := path.Join(home, ".kube", "kube-common")
  return dirPath, nil
}
