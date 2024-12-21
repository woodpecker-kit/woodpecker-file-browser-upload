[![ci](https://github.com/woodpecker-kit/woodpecker-file-browser-upload/workflows/ci/badge.svg)](https://github.com/woodpecker-kit/woodpecker-file-browser-upload/actions/workflows/ci.yml)

[![go mod version](https://img.shields.io/github/go-mod/go-version/woodpecker-kit/woodpecker-file-browser-upload?label=go.mod)](https://github.com/woodpecker-kit/woodpecker-file-browser-upload)
[![GoDoc](https://godoc.org/github.com/woodpecker-kit/woodpecker-file-browser-upload?status.png)](https://godoc.org/github.com/woodpecker-kit/woodpecker-file-browser-upload)
[![goreportcard](https://goreportcard.com/badge/github.com/woodpecker-kit/woodpecker-file-browser-upload)](https://goreportcard.com/report/github.com/woodpecker-kit/woodpecker-file-browser-upload)

[![GitHub license](https://img.shields.io/github/license/woodpecker-kit/woodpecker-file-browser-upload)](https://github.com/woodpecker-kit/woodpecker-file-browser-upload)
[![codecov](https://codecov.io/gh/woodpecker-kit/woodpecker-file-browser-upload/branch/main/graph/badge.svg)](https://codecov.io/gh/woodpecker-kit/woodpecker-file-browser-upload)
[![GitHub latest SemVer tag)](https://img.shields.io/github/v/tag/woodpecker-kit/woodpecker-file-browser-upload)](https://github.com/woodpecker-kit/woodpecker-file-browser-upload/tags)
[![GitHub release)](https://img.shields.io/github/v/release/woodpecker-kit/woodpecker-file-browser-upload)](https://github.com/woodpecker-kit/woodpecker-file-browser-upload/releases)

## for what

- this project used to woodpecker CI to support [file browser](https://github.com/filebrowser/filebrowser)
- get file browser `host` like `https://filebrowser.xxx.com`
- this plugin need file browser `username`
- and need file browser user need `password`

## Contributing

[![Contributor Covenant](https://img.shields.io/badge/contributor%20covenant-v1.4-ff69b4.svg)](.github/CONTRIBUTING_DOC/CODE_OF_CONDUCT.md)
[![GitHub contributors](https://img.shields.io/github/contributors/woodpecker-kit/woodpecker-file-browser-upload)](https://github.com/woodpecker-kit/woodpecker-file-browser-upload/graphs/contributors)

We welcome community contributions to this project.

Please read [Contributor Guide](.github/CONTRIBUTING_DOC/CONTRIBUTING.md) for more information on how to get started.

请阅读有关 [贡献者指南](.github/CONTRIBUTING_DOC/zh-CN/CONTRIBUTING.md) 以获取更多如何入门的信息

## Features

- [x] upload file to file browser
    - [x] support file with [glob](https://pkg.go.dev/path/filepath#Match)
    - [x] support file with regular expression
- [x] set file custom dist graph
- [x] support share link
    - [x] support share link expires
    - [x] support share link auto password
    - [x] support share link password
- [x] support send result to woodpecker steps transfer
- [x] docker platform support (v1.6.+)
  -  `linux/amd64 linux/386 linux/arm64/v8 linux/arm/v7 linux/ppc64le linux/s390x`
- [x] support link choose by web test
  - [x] `file-browser-urls` support multi urls, will auto switch host fast (1.7.+)

## usage

### workflow usage

- see [doc](doc/docs.md)

---

- want dev this project, see [doc](doc/README.md)