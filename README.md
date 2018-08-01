# kubectl-common

A wrapper around `kubectl` to easily switch between kubectl versions and
kubeconfigs.

`kubectl-common` uses aliases to refer to configurations that provide the
required kubectl version and kubeconfig path.

This tool is useful for:
* **Switching between clusters that require different kubeconfig files.**
Changing the kubeconfig path can be tedious, annoying when dealing with multiple
terminals (if using env variable KUBECONFIG), and verbose if explicitly using the
--kubeconfig flag with kubectl.
* **Switching between clusters that require different versions of kubectl.**
Using different kubectl command names for different versions is annoying and
easy to mess up.

## Installing

Download kubectl-common from [the releases page.](https://github.com/tkporter/kubectl-common/releases)

Once you've expanded the `.tar.gz`, move `kubectl-common` to somewhere in your
$PATH. The easiest way to use this tool is to move it to `/usr/local/bin/` and
name it `kubectl` so it can be used exactly the same as an original `kubectl`.

## Alias Configuration

Aliases are defined in an `alias-config.json` file. Here's an example of what
it would look like:

```json
{
  "aliases": {
    "foo": {
      "version": "1.8.4",
      "kubeconfig": "/Users/johnsmith/.kube/config_foo"
    },
    "bar": {
      "version": "1.10.0",
      "kubeconfig": "/Users/johnsmith/.kube/config_bar"
    },
    "baz": {
      "version": "1.11.1",
      "kubeconfig": "/Users/johnsmith/.kube/config_baz"
    },
    "zoo": {
      "version": "1.11.1",
      "kubeconfig": "/Users/johnsmith/.kube/config_zoo"
    }
  }
}
```

Where the alias "foo" would correspond to `kubectl` version `1.8.4` and
kubeconfig `/Users/johnsmith/.kube/config_foo`, "bar" to version `1.10.0` and
kubeconfig `/Users/johnsmith/.kube/config_bar`, etc. It's fine to have more
than one alias with the same version.

Note that full absolute paths are required for the `kubeconfig` (`~` is not
supported).

## Usage

For the easiest experience, it's recommended to rename `kubectl-common` to
`kubectl`.

### Applying an Alias Config

To apply an `alias-config.json`, run:

```
$ kubectl-common apply-alias-config
```

By default, this looks for an `alias-config.json` in `~/.kube/kube-common`.
If your `alias-config.json` is in a different directory, you can specify
the directory using the `-d` or `--dir` flag. Once this command is run,
you won't need to re-specify the directory of the alias config unless if
running `kubectl-common apply-alias-config` again in the future.

### Switching Aliases

We can change our current kubectl version/kubeconfig configuration by:

```
$ kubectl-common use-alias foo
Now using alias foo for kubectl version 1.8.4 and kubeconfig /Users/johnsmith/.kube/config_foo
```

### Performing Normal kubectl Commands

Once an alias has been selected using `use-alias`, we can just use
`kubectl-common` as if it were `kubectl` directly:

```
$ kubectl-common version
# verify we are using the right kubectl version!
$ kubectl-common get pods
# stuff
$ kubectl-common get deployments
# things
```

### Next Steps

Short term goals:
* Add tests
* Support & test on platforms other than macOS
* Allow `~` usage in the alias config for kubeconfigs. This has to do with
the way we're executing the actual kubectl command because bash is required to
replacing `~` with `$HOME`.
* Add flag to `apply-alias-config` to determine where to store kubectl versions
that are copied or downloaded.
