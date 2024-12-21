# woodpecker-file-browser-upload

how to dev

## env

- minimum go version: go 1.21
- change `go 1.21`, `^1.21`, `1.21.13` to new go version

### libs

| lib                                                | version |
|:---------------------------------------------------|:--------|
| https://github.com/stretchr/testify                | v1.8.4  |
| https://github.com/sebdah/goldie                   | v2.5.3  |
| https://github.com/gookit/color                    | v1.5.3  |
| https://github.com/urfave/cli/                     | v2.23.7 |
| https://github.com/woodpecker-kit/woodpecker-tools | v1.5.0  |

- more libs see [go.mod](https://github.com/woodpecker-kit/woodpecker-file-browser-upload/blob/main/go.mod)

## depends

in go mod project

```bash
# warning use private git host must set
# global set for once
# add private git host like github.com to evn GOPRIVATE
$ go env -w GOPRIVATE='github.com'
# use ssh proxy
# set ssh-key to use ssh as http
$ git config --global url."git@github.com:".insteadOf "https://github.com/"
# or use PRIVATE-TOKEN
# set PRIVATE-TOKEN as gitlab or gitea
$ git config --global http.extraheader "PRIVATE-TOKEN: {PRIVATE-TOKEN}"
# set this rep to download ssh as https use PRIVATE-TOKEN
$ git config --global url."ssh://github.com/".insteadOf "https://github.com/"

# before above global settings
# test version info
$ git ls-remote -q https://github.com/woodpecker-kit/woodpecker-file-browser-upload.git

# test depends see full version
$ go list -mod readonly -v -m -versions github.com/woodpecker-kit/woodpecker-file-browser-upload
# or use last version add go.mod by script
$ echo "go mod edit -require=$(go list -mod=readonly -m -versions github.com/woodpecker-kit/woodpecker-file-browser-upload | awk '{print $1 "@" $NF}')"
$ echo "go mod vendor"
```

## local dev

```bash
# It needs to be executed after the first use or update of dependencies.
$ make init dep
```

- test code

```bash
$ make test testBenchmark
```

add main.go file and run

```bash
# run and shell help
$ make devHelp

# run at PLUGIN_DEBUG=true
$ make dev

# run at ordinary mode
$ make run
```

- ci to fast check

```bash
# check style at local
$ make style

# run ci at local
$ make ci
```

### docker

```bash
# then test build as test/Dockerfile
$ make dockerTestRestartLatest
# clean test build
$ make dockerTestPruneLatest

# more info see
$ make helpDocker
```

### EngineeringStructure

```
.
├── Dockerfile                     # ci docker build
├── build.dockerfile               # local docker build
├── Makefile                       # make entry
├── README.md
├── build                          # build output folder
├── dist                           # dist output folder
├── cmd
│     ├── cli
│     │     ├── app.go             # cli entry
│     │     ├── cli_aciton_test.go # cli action test
│     │     └── cli_action.go      # cli action
│     └── woodpecker-plugin-env    # command line main package install and dev entrance
│         ├── main.go                   # command line entry
│         └── main_test.go              # integrated test entry
├── constant                       # constant package
│         ├── common_flag.go         # common environment variable
│         ├── platform_flag.go       # platform environment variable
│         └── version.go             # semver version constraint set
├── doc
│         ├── README.md              # command line tools documentation
│         └── docs.md                # woodpecker documentation
├── go.mod
├── go.sum
├── z-MakefileUtils                # make toolkit
├── package.json                   # command line profile information for embed
├── resource.go                    # embed resource
├── internal                          # toolkit package
│         ├── pkgJson                 # package.json toolkit
│         └── version_check           # version check by semver
├── plugin                         # plugin package
│         ├── flag.go                 # plugin flag
│         ├── impl.go                 # plugin implement
│         ├── plugin.go               # plugin entry
│         └── settings.go             # plugin settings
├── plugin_test                    # plugin test
│         ├── init_test.go            # each test init
│         └── plugin_test.go          # plugin test
└── zymosis                         # resource mark by https://github.com/convention-change/zymosis

```

### log

- open debug log by env `PLUGIN_DEBUG=true` or global flag `--plugin.debug true`

```go
package foo

func GlobalBeforeAction(c *cli.Context) error {
  isDebug := wd_urfave_cli_v2.IsBuildDebugOpen(c)
  if isDebug {
    wd_log.OpenDebug()
  }
  return nil
}
```

### template

- [https://github.com/aymerick/raymond](https://github.com/aymerick/raymond)
- function doc [https://masterminds.github.io/sprig/](https://masterminds.github.io/sprig/)

- open template support at cli `main.go`

```go
package main
func main() {
	// register helpers once
	wd_template.RegisterSettings(wd_template.DefaultHelpers)
}
```
- and open at test `init_test.go`

```go
package plugin_test

func init() {
	// if open wd_template please open this
	wd_template.RegisterSettings(wd_template.DefaultHelpers)
}
```