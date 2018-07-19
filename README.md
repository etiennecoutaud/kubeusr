# Kubeusr
CLI tools to easily manage kubernetes user based on service account and token connection.


`kubeusr` can create to kind of user:
* cluster wide admin
* namespace admin whith scoped rights in the namespace

[![asciicast](https://asciinema.org/a/fgAtKRyhOTOYA2zSvIiRZWXbm.png)](https://asciinema.org/a/fgAtKRyhOTOYA2zSvIiRZWXbm?speed=1.7)

## Install

```
git clone http://github.com/etiennecoutaud/kubeusr
cd kubeusr
make install
```

## Usage
```
$ kubeusr
kubeusr manage Kubernetes users

Usage:
  kubeusr [command]

Available Commands:
  create      Create Kubernetes user
  delete      Delete Kubernetes user
  get         Get or List Kubernetes user
  help        Help about any command
  kubeconfig  Get user kubeconfig

Flags:
  -h, --help                help for kubeusr
      --kubeconfig string   absolute path to the kubeconfig (default is KUBECONFIG env var value)
  -t, --toggle              Help message for toggle

Use "kubeusr [command] --help" for more information about a command.
```
