package kubectl_manager

import (
  "fmt"
  "io"
  "net/http"
  "os"
  "os/exec"
  "path"
  "regexp"
  "syscall"
)

const baseBinPath = "/usr/local/bin/"
var kubectlPath = path.Join(baseBinPath, "kubectl")
var executableFileMode = os.FileMode(0755)

// Executes a command defined by args using kubectl version `version`
func RunKubectlCommand(version, kubeconfig string, args []string) {
  kubectlPath := getVersionedKubectlPath(version)
  fullArgs := append([]string{ kubectlPath, "--kubeconfig", kubeconfig }, args...)

  err := syscall.Exec(kubectlPath, fullArgs, os.Environ())
  if err != nil {
    panic(err)
  }
}

// Ensures there is a local copy of kubectl for each version in the map
func SetupKubectlVersions(versions map[string]bool) {
  existingVersion := getExistingKubectlVersion()
  for version := range versions {
    setupKubectlVersion(version, existingVersion)
  }
}

// Gets the entire path of a kubectl for a particular version
func getVersionedKubectlPath(version string) string {
  return fmt.Sprintf("%s_v%s", kubectlPath, version)
}

// Downloads the kubectl for a particular version if not found locally.
// If the kubectl for a version is found locally, it's copied to the desired
// path.
// TODO(tkporter): check the versioned kubectl path as well. Better error handling?
func setupKubectlVersion(version, existingVersion string) {
  newKubeVersionPath := getVersionedKubectlPath(version)

  // If we already have one of the existing versions, just rename it
  if existingVersion == version {
    fmt.Printf("kubectl v%s found at %s\nCopying existing file to %s\n\n", version, kubectlPath, newKubeVersionPath)
    copyFile(kubectlPath, newKubeVersionPath)
  } else {
    fmt.Printf("Downloading kubectl v%s to %s\n\n", version, newKubeVersionPath)
    err := downloadFile(newKubeVersionPath, getKubectlDownloadUrl(version))
    if err != nil {
      panic(err)
    }
  }
  // Ensure the kubectl is executable
  err := os.Chmod(newKubeVersionPath, executableFileMode)
  if err != nil {
    panic(err)
  }
}

// Gets the version of an existing kubectl at kubectlPath
// TODO(tkporter): better error handling? Maybe extend to versioned paths too?
func getExistingKubectlVersion() string {
  versionOutput, err := exec.Command(kubectlPath, "version", "--client", "--short").Output()
  if err != nil {
    panic(err)
    return ""
  }
  re := regexp.MustCompile("\\d+\\.\\d+\\.\\d+")
  existingVersion := re.FindString(string(versionOutput))
  return existingVersion
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
