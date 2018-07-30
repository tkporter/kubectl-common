package kubectl_manager

import (
  "fmt"
  "io"
  "net/http"
  "os"
  "os/exec"
  "path"
  "path/filepath"
  "regexp"
  "strings"
  "syscall"
)

const baseBinPath = "/usr/local/bin/"
var kubectlPath = path.Join(baseBinPath, "kubectl")
var executableFileMode = os.FileMode(0755)

// Executes a command defined by args using kubectl version `version`
func RunKubectlCommand(version, kubeconfig string, args []string) error {
  kubectlPath := getVersionedKubectlPath(version)
  fullArgs := append([]string{ kubectlPath, "--kubeconfig", kubeconfig }, args...)

  // returns error
  return syscall.Exec(kubectlPath, fullArgs, os.Environ())
}

// Ensures there is a local copy of kubectl for each version in the map
func SetupKubectlVersions(versions map[string]bool) error {
  existingVersions, err := getExistingKubectlVersions()
  if err != nil {
    return err
  }
  for version := range versions {
    err := setupKubectlVersion(version, existingVersions)
    if err != nil {
      return err
    }
  }
  return nil
}

// Gets the entire path of a kubectl for a particular version
func getVersionedKubectlPath(version string) string {
  return fmt.Sprintf("%s_v%s", kubectlPath, version)
}

// Downloads the kubectl for a particular version if not found locally.
// If the kubectl for a version is found locally, it's copied to the desired
// path.
func setupKubectlVersion(version string, existingVersions map[string]string) error {
  newKubeVersionPath := getVersionedKubectlPath(version)

  // If we already have one of the existing versions, rename it if the path
  // is bad (not in proper kubectl_vX.X.X format) or just continue if the path
  // is already correct
  if kcPath, ok := existingVersions[version]; ok {
    fmt.Printf("kubectl v%s found at %s\n", version, kcPath)
    if kcPath == newKubeVersionPath {
      fmt.Printf("Path is proper version format, continuing.\n\n")
    } else {
      fmt.Printf("Copying existing file to %s\n\n", newKubeVersionPath)
      copyFile(kubectlPath, newKubeVersionPath)
    }
  } else {
    fmt.Printf("Downloading kubectl v%s to %s\n\n", version, newKubeVersionPath)
    err := downloadFile(newKubeVersionPath, getKubectlDownloadUrl(version))
    if err != nil {
      return err
    }
  }
  // Ensure the kubectl is executable, returns err
  return os.Chmod(newKubeVersionPath, executableFileMode)
}

// Gets all existing versions of the form `${kubectlPath}_v*`, and returns
// a map with versions as keys and full paths of the files as values
func getExistingKubectlVersions() (map[string]string, error) {
  matches, err := filepath.Glob(strings.Join([]string{ kubectlPath, "_v*" }, ""))
  if err != nil {
    return nil, err
  }
  versionMap := make(map[string]string)
  _ = versionMap
  for _, kcPath := range matches {
    version, err := getKubectlVersion(kcPath)
    if err != nil {
      // Error getting version of existing kubectl, just ignore because
      // we will just redownload if this version is needed by our config file
      continue
    }
    versionMap[version] = kcPath
  }
  return versionMap, nil
}

// Gets the version of an existing kubectl at kcPath
// TODO(tkporter): better error handling?
func getKubectlVersion(kcPath string) (string, error) {
  versionOutput, err := exec.Command(kcPath, "version", "--client", "--short").Output()
  if err != nil {
    return "", err
  }
  re := regexp.MustCompile("\\d+\\.\\d+\\.\\d+")
  existingVersion := re.FindString(string(versionOutput))
  if existingVersion == "" {
    return "", fmt.Errorf("No valid version found")
  }
  return existingVersion, nil
}

// Returns the url to download a particular version of kubectl
// TODO(tkporter): extend this to other platforms
func getKubectlDownloadUrl(version string) string {
  // just for darwin for now
  return fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/v%s/bin/darwin/amd64/kubectl", version)
}

// Copies a file at path src to dst
func copyFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, in)
    if err != nil {
        return err
    }
    return out.Close()
}

// Downloads a file at url to filepath
func downloadFile(filepath, url string) error {
  // Create the file
  out, err := os.Create(filepath)
  if err != nil  {
    return err
  }
  defer out.Close()

  // Get the data
  resp, err := http.Get(url)
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  // Check server response
  if resp.StatusCode != http.StatusOK {
    return fmt.Errorf("bad status: %s", resp.Status)
  }

  // Writer the body to file
  _, err = io.Copy(out, resp.Body)
  if err != nil  {
    return err
  }
  return nil
}
