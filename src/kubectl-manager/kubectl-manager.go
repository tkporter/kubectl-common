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
var binFileMode = os.FileMode(0755)

func RunKubectlCommand(version string, args []string) {
  kubectlPath := getVersionedKubectlPath(version)
  argsWithKubectl := append([]string{ kubectlPath }, args...)

  err := syscall.Exec(kubectlPath, argsWithKubectl, os.Environ())
  if err != nil {
    panic(err)
  }
}

func SetupKubectlVersions(versions map[string]bool) {
  existingVersion := getExistingKubectlVersion()
  for version := range versions {
    setupKubectlVersion(version, existingVersion)
  }
}

func getVersionedKubectlPath(version string) string {
  return fmt.Sprintf("%s_v%s", kubectlPath, version)
}

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

  err := os.Chmod(newKubeVersionPath, binFileMode)
  if err != nil {
    panic(err)
  }
}

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

func getKubectlDownloadUrl(version string) string {
  // just for macs for now
  return fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/v%s/bin/darwin/amd64/kubectl", version)
}

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
