package file_browser_upload

import (
	"github.com/woodpecker-kit/woodpecker-tools/wd_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_short_info"
	"github.com/woodpecker-kit/woodpecker-transfer-data/wd_share_file_browser_upload"
)

type (
	// FileBrowserPlugin file_browser_upload all config
	FileBrowserPlugin struct {
		Name           string
		Version        string
		woodpeckerInfo *wd_info.WoodpeckerInfo
		wdShortInfo    *wd_short_info.WoodpeckerInfoShort
		onlyArgsCheck  bool

		Settings Settings

		FuncPlugin FuncPlugin `json:"-"`

		shareFileBrowserUpload *wd_share_file_browser_upload.WdShareFileBrowserUpload
	}
)

type FuncPlugin interface {
	ShortInfo() wd_short_info.WoodpeckerInfoShort

	SetWoodpeckerInfo(info wd_info.WoodpeckerInfo)
	GetWoodPeckerInfo() wd_info.WoodpeckerInfo

	OnlyArgsCheck()

	Exec() error

	loadStepsTransfer() error
	checkArgs() error
	saveStepsTransfer() error
}
