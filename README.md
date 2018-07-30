# kubectl-common

A wrapper around `kubectl` to easily switch between versions.

## Alias Configuration

`kubectl-common` uses aliases of versions to easily switch between versions
without needing to memorize the numeric `x.x.x` version of `kubectl` needed.

Aliases are defined in an `alias-config.json` file. Here's an example of what
it would look like:

```json
{
  "aliases": {
    "foo": "1.8.4",
    "bar": "1.10.0",
    "baz": "1.11.1",
    "zoo": "1.11.1"
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

### Switching Version Aliases

We can change our current kubectl version to use by:

```
$ kubectl-common use-version-alias foo
Now using alias foo for kubectl version 1.8.4
```

### Performing Normal kubectl Commands

Once an alias has been selected using `use-version-alias`, we can just use
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

This will probably be fixed soon, but:

`kubectl-common --help` will show the help for `kubectl-common`, rather than
that of `kubectl`.

If you want to view `kubectl`'s help message, use `kubectl-common -- --help`.

This includes if you want to view a help message for a command in `kubectl`:
`kubectl-common -- get --help`.
