# Kubeusr
CLI tools to easily manage kubernetes user based on service account and token connection.


`kubeusr` can create two kind of user:
* cluster wide admin
* namespace admin whith scoped rights in the namespace

[![asciicast](https://asciinema.org/a/fgAtKRyhOTOYA2zSvIiRZWXbm.png)](https://asciinema.org/a/fgAtKRyhOTOYA2zSvIiRZWXbm?speed=1.7)

## Install

```
$ mkdir -p $GOPATH/src/github.com/etiennecoutaud
$ cd $GOPATH/src/github.com/etiennecoutaud
$ git clone http://github.com/etiennecoutaud/kubeusr
$ cd kubeusr
$ make install
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

## Create multiple user from file

File format:
```
$ cat list.yml
---
list:
- name: user1
  scope: namespace
  namespace: foo
- name: user2
  scope: namespace
  namespace: foo
- name: user3
  scope: cluster
  namespace: default
```

```
$ kubeusr create -f list.yml
User "user1" created with scoped rights in namespace "foo"
User "user2" created with scoped rights in namespace "foo"
User "user3" created with cluster-wide rights in namespace "default"

$ kubeusr get
NAMESPACE   USER    SCOPE       CREATION DATE
default     user3   cluster     2018-07-19 17:23:57 +0200 CEST
foo         user1   namespace   2018-07-19 17:23:57 +0200 CEST
foo         user2   namespace   2018-07-19 17:23:57 +0200 CEST
```
