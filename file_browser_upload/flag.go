package file_browser_upload

import (
	"fmt"
	"github.com/sinlov/filebrowser-client/web_api"
	"github.com/urfave/cli/v2"
	"github.com/woodpecker-kit/woodpecker-tools/wd_flag"
	"github.com/woodpecker-kit/woodpecker-tools/wd_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_log"
)

const (
	CliNameFileBrowserHost = "settings.file_browser_host"
	EnvFileBrowserHost     = "PLUGIN_FILE_BROWSER_HOST"

	CliNameFileBrowserUsername = "settings.file_browser_username"
	EnvFileBrowserUsername     = "PLUGIN_FILE_BROWSER_USERNAME"

	CliNameFileBrowserUserPassword = "settings.file_browser_user_password"
	EnvFileBrowserUserPassword     = "PLUGIN_FILE_BROWSER_USER_PASSWORD"

	CliNameFileBrowserTimeOutSendSecond = "settings.file_browser_timeout_send_second"
	EnvFileBrowserTimeOutSendSecond     = "PLUGIN_FILE_BROWSER_TIMEOUT_SEND_SECOND"

	CilNameFileBrowserWorkSpace = "settings.file_browser_work_space"
	EnvFileBrowserWorkSpace     = "PLUGIN_FILE_BROWSER_WORK_SPACE"

	CliNameFileBrowserDistType = "settings.file_browser_dist_type"
	EnvFileBrowserDistType     = "PLUGIN_FILE_BROWSER_DIST_TYPE"

	CliNameFileBrowserDistGraph = "settings.file_browser_dist_graph"
	EnvFileBrowserDistGraph     = "PLUGIN_FILE_BROWSER_DIST_GRAPH"

	CliNameFileBrowserRemoteRootPath = "settings.file_browser_remote_root_path"
	EnvFileBrowserRemoteRootPath     = "PLUGIN_FILE_BROWSER_REMOTE_ROOT_PATH"

	CliNameFileBrowserTargetDistRootPath = "settings.file_browser_target_dist_root_path"
	EnvFileBrowserTargetDistRootPath     = "PLUGIN_FILE_BROWSER_TARGET_DIST_ROOT_PATH"

	CliNameFileBrowserFileGlob = "settings.file_browser_file_glob"
	EnvFileBrowserFileGlob     = "PLUGIN_FILE_BROWSER_FILE_GLOB"

	CliNameFileBrowserFileRegular = "settings.file_browser_file_regular"
	EnvFileBrowserFileRegular     = "PLUGIN_FILE_BROWSER_FILE_REGULAR"

	CliNameFileBrowserShareLinkEnable = "settings.file_browser_share_link_enable"
	EnvFileBrowserShareLinkEnable     = "PLUGIN_FILE_BROWSER_SHARE_LINK_ENABLE"

	CilNameFileBrowserShareLinkUnit = "settings.file_browser_share_link_unit"
	EnvFileBrowserShareLinkUnit     = "PLUGIN_FILE_BROWSER_SHARE_LINK_UNIT"

	CliNameFileBrowserShareLinkExpire = "settings.file_browser_share_link_expire"
	EnvFileBrowserShareLinkExpire     = "PLUGIN_FILE_BROWSER_SHARE_LINK_EXPIRE"

	CliNameFileBrowserShareAutoPasswordEnable = "settings.file_browser_share_auto_password_enable"
	EnvFileBrowserShareAutoPasswordEnable     = "PLUGIN_FILE_BROWSER_SHARE_AUTO_PASSWORD_ENABLE"

	CliNameFileBrowserShareLinkPasswd = "settings.file_browser_share_link_passwd"
	EnvFileBrowserShareLinkPasswd     = "PLUGIN_FILE_BROWSER_SHARE_LINK_PASSWD"
)

// GlobalFlag
// Other modules also have flags
func GlobalFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    CliNameFileBrowserHost,
			Usage:   "file browser host",
			EnvVars: []string{EnvFileBrowserHost},
		},
		&cli.StringFlag{
			Name:    CliNameFileBrowserUsername,
			Usage:   "file browser username",
			EnvVars: []string{EnvFileBrowserUsername},
		},
		&cli.StringFlag{
			Name:    CliNameFileBrowserUserPassword,
			Usage:   "file browser user password",
			EnvVars: []string{EnvFileBrowserUserPassword},
		},
		&cli.UintFlag{
			Name:    CliNameFileBrowserTimeOutSendSecond,
			Usage:   "file browser timeout send second",
			EnvVars: []string{EnvFileBrowserTimeOutSendSecond},
		},
		&cli.StringFlag{
			Name:    CilNameFileBrowserWorkSpace,
			Usage:   fmt.Sprintf("file_browser work space. default will use env: %s", wd_flag.EnvKeyNameCiWorkspace),
			EnvVars: []string{EnvFileBrowserWorkSpace},
		},
		&cli.StringFlag{
			Name:    CliNameFileBrowserDistType,
			Usage:   fmt.Sprintf("file browser dist type, only support %v", pluginDistTypeSupport),
			EnvVars: []string{EnvFileBrowserDistType},
		},
		&cli.StringFlag{
			Name:    CliNameFileBrowserDistGraph,
			Usage:   "file browser dist graph",
			EnvVars: []string{EnvFileBrowserDistGraph},
			Value:   DistGraphTypeGit,
		},
		&cli.StringFlag{
			Name:    CliNameFileBrowserRemoteRootPath,
			Usage:   "must set args, this will append by file_browser_dist_type at remote",
			EnvVars: []string{EnvFileBrowserRemoteRootPath},
		},
		&cli.StringFlag{
			Name:    CliNameFileBrowserTargetDistRootPath,
			Usage:   "path of file_browser local work on root, can set \"\"",
			Value:   "",
			EnvVars: []string{EnvFileBrowserTargetDistRootPath},
		},
		&cli.StringSliceFlag{
			Name:    CliNameFileBrowserFileGlob,
			Usage:   "must set args, globs list of send to file_browser under file_browser_target_dist_root_path",
			EnvVars: []string{EnvFileBrowserFileGlob},
		},
		&cli.StringFlag{
			Name:    CliNameFileBrowserFileRegular,
			Usage:   "must set args, regular of send to file_browser under file_browser_target_dist_root_path",
			EnvVars: []string{EnvFileBrowserFileRegular},
		},
		&cli.BoolFlag{
			Name:    CliNameFileBrowserShareLinkEnable,
			Usage:   "file browser share link enable",
			EnvVars: []string{EnvFileBrowserShareLinkEnable},
		},
		&cli.StringFlag{
			Name:    CilNameFileBrowserShareLinkUnit,
			Usage:   fmt.Sprintf("file browser share link unit, only support %v", web_api.ShareUnitDefine()),
			EnvVars: []string{EnvFileBrowserShareLinkUnit},
			Value:   web_api.ShareUnitDays,
		},
		&cli.UintFlag{
			Name:    CliNameFileBrowserShareLinkExpire,
			Usage:   "if set 0, will allow share_link exist forever, default: 0",
			Value:   0,
			EnvVars: []string{EnvFileBrowserShareLinkExpire},
		},
		&cli.BoolFlag{
			Name:    CliNameFileBrowserShareAutoPasswordEnable,
			Usage:   "password of share_link auto , if open this will cover settings.file_browser_share_link_password",
			Value:   false,
			EnvVars: []string{EnvFileBrowserShareAutoPasswordEnable},
		},
		&cli.StringFlag{
			Name:    CliNameFileBrowserShareLinkPasswd,
			Usage:   "password of share_link, if not set will not use password, default: \"\"",
			Value:   "",
			EnvVars: []string{EnvFileBrowserShareLinkPasswd},
		},
	}
}

func HideGlobalFlag() []cli.Flag {
	return []cli.Flag{}
}

func BindCliFlags(c *cli.Context,
	debug bool,
	cliName, cliVersion string,
	wdInfo *wd_info.WoodpeckerInfo, rootPath,
	stepsTransferPath string, stepsOutDisable bool,
) (*FileBrowserPlugin, error) {

	config := Settings{
		Debug:             debug,
		TimeoutSecond:     c.Uint(wd_flag.NameCliPluginTimeoutSecond),
		StepsTransferPath: stepsTransferPath,
		StepsOutDisable:   stepsOutDisable,
		RootPath:          rootPath,

		FileBrowserBaseConfig: FileBrowserBaseConfig{
			FileBrowserHost:              c.String(CliNameFileBrowserHost),
			FileBrowserUsername:          c.String(CliNameFileBrowserUsername),
			FileBrowserUserPassword:      c.String(CliNameFileBrowserUserPassword),
			FileBrowserTimeoutPushSecond: c.Uint(CliNameFileBrowserTimeOutSendSecond),
			FileBrowserWorkSpace:         c.String(CilNameFileBrowserWorkSpace),
		},

		FileBrowserSendConfig: FileBrowserSendConfig{
			FileBrowserRemoteRootPath:              c.String(CliNameFileBrowserRemoteRootPath),
			FileBrowserDistType:                    c.String(CliNameFileBrowserDistType),
			FileBrowserDistGraph:                   c.String(CliNameFileBrowserDistGraph),
			FileBrowserTargetDistRootPath:          c.String(CliNameFileBrowserTargetDistRootPath),
			FileBrowserTargetFileGlob:              c.StringSlice(CliNameFileBrowserFileGlob),
			FileBrowserTargetFileRegular:           c.String(CliNameFileBrowserFileRegular),
			FileBrowserShareLinkEnable:             c.Bool(CliNameFileBrowserShareLinkEnable),
			FileBrowserShareLinkUnit:               c.String(CilNameFileBrowserShareLinkUnit),
			FileBrowserShareLinkExpires:            c.Uint(CliNameFileBrowserShareLinkExpire),
			FileBrowserShareLinkAutoPasswordEnable: c.Bool(CliNameFileBrowserShareAutoPasswordEnable),
			FileBrowserShareLinkPassword:           c.String(CliNameFileBrowserShareLinkPasswd),
		},
	}

	// set default TimeoutSecond
	if config.TimeoutSecond == 0 {
		config.TimeoutSecond = 10
	}

	wd_log.Debugf("args %s: %v", wd_flag.NameCliPluginTimeoutSecond, config.TimeoutSecond)

	p := FileBrowserPlugin{
		Name:           cliName,
		Version:        cliVersion,
		woodpeckerInfo: wdInfo,
		Settings:       config,
	}

	return &p, nil
}
