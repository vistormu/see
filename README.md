<a name="readme-top"></a>

<div align="center">

<a href="https://github.com/vistormu/see" target="_blank" title="go to the repo"><img width="196px" alt="see logo" src="/docs/logo.svg"></a>


# see
# a better way to visualize your file system

_see_ is the replacement of `ls`, `tree`, and `cat` commands with a more user-friendly output, with a focus on `git` repositories.

[![go version][go_version_img]][go_dev_url]
[![Go report][go_report_img]][go_report_url]
[![License][repo_license_img]][repo_license_url]

<a href="https://github.com/vistormu/see" target="_blank" title=""><img width="99%" alt="see command" src="/docs/ls.png"></a>

<a href="https://github.com/vistormu/see" target="_blank" title=""><img width="99%" alt="see command" src="/docs/cat.png"></a>

</div>

> [!WARNING]
> this project is functional but still in development, so expect some bugs and missing features

## ✨ features

- colorful and pretty output
- see directly the status of your git repositories:
    - green means the repository is clean
    - yellow means there are uncommitted changes
    - red means there are uncommitted changes and untracked files
- if `zoxide` is installed, WIP

## ⚡️ quick start

just type

```bash
see
```

to see the current directory content, or

```bash
see <path>
```

to see the content of a specific path or file.


### flags

| flag | description | status |
| --- | --- | --- |
| `-h`, `--help` | show help | ❌ |
| `-v`, `--version` | show version | ❌ |
| `-f`, `--filter` | filter the output by a specific string (e.g. `see -f .txt`) | ❌ |
| `-d`, `--depth` | set the depth of the tree (default: 1) | ❌ |
| `-s`, `--sort` | sort files by name, kind, size, git status, or date (default: name) | ❌ |
| `-n`, `--nerd` | show all possible information about the tree | ❌ |

<!-- ✅ -->

## 🚀 installation

### homebrew

> [!WARNING]
> not implemented

if you have [homebrew](https://brew.sh/) installed, you can tap the formula

```bash
brew tap vistormu/see
```

and then install _see_ with the following command:

```bash
brew install vistormu/see
```

### using go

if you have `go` installed, you can install _see_ with the following command:

```bash
go install github.com/vistormu/see@latest
```

this will install the binary in your `$GOPATH/bin` directory, so make sure to add it to your `PATH` if it's not already there.

### from source

clone the repo:

```bash
git clone https://github.com/vistormu/see.git
```

then build the project:

```bash
cd see
go build -o see
```

(optional) move the binary to a directory in your `PATH`:

```bash
mv see /usr/local/bin/
```

## 🌟 stargazers

[![Stargazers over time](https://starchart.cc/vistormu/see.svg?variant=adaptive)](https://starchart.cc/vistormu/see)

<div align="right">

[&nwarr; Back to top](#readme-top)

</div>


<!-- links -->
[go_version_img]: https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go
[go_dev_url]: https://go.dev/
[go_report_img]: https://goreportcard.com/badge/github.com/vistormu/see
[go_report_url]: https://goreportcard.com/report/github.com/vistormu/see
[repo_license_img]: https://img.shields.io/github/license/vistormu/see?style=for-the-badge
[repo_license_url]: https://github.com/vistormu/see/blob/main/LICENSE
