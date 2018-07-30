# kubectl-common

A wrapper around `kubectl` to easily switch between kubectl versions and
kubeconfigs.

`kubectl-common` uses aliases to refer to configurations that provide the
required kubectl version and kubeconfig path.

This tool is useful for:
* Switching between clusters that require different kubeconfig files. Explicitly
changing the kubeconfig path can be tedious.
* Switching between clusters that require different versions of kubectl.
Using different kubectl command names for different versions is annoying and
easy to mess up.

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
    },
  }
}
```

Where the alias "foo" would correspond to `kubectl` version `1.8.4`, "bar" to
version `1.10.0`, etc. It's fine to have more than one alias with the same
version.

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

### Quirks

`kubectl-common --help` will show the help for `kubectl-common`, rather than
that of `kubectl`.

If you want to view `kubectl`'s help message, use `kubectl-common -- --help`.

This includes if you want to view a help message for a command in `kubectl`:
`kubectl-common -- get --help`.

### Next Steps

Short term goals:
* Support & test on platforms other than macOS
* Allow `~` usage in the alias config for kubeconfigs
