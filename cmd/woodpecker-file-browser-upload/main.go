//go:build !test

package main

import (
	"github.com/gookit/color"
	"github.com/woodpecker-kit/woodpecker-file-browser-upload"
	"github.com/woodpecker-kit/woodpecker-file-browser-upload/cmd/cli"
	"github.com/woodpecker-kit/woodpecker-file-browser-upload/internal/pkgJson"
	"github.com/woodpecker-kit/woodpecker-tools/wd_log"
	"github.com/woodpecker-kit/woodpecker-tools/wd_template"
	os "os"
)

func main() {
	wd_log.SetLogLineDeep(wd_log.DefaultExtLogLineMaxDeep)
	pkgJson.InitPkgJsonContent(woodpecker_file_browser_upload.PackageJson)

	// register helpers once
	wd_template.RegisterSettings(wd_template.DefaultHelpers)

	app := cli.NewCliApp()

	args := os.Args
	if err := app.Run(args); nil != err {
		color.Redf("cli err at %v\n", err)
	}
}
