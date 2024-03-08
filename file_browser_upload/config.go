package file_browser_upload

const (
	DistTypeGit    = "git"
	DistTypeCustom = "custom"

	DistGraphTypeGit = "{{ Repo.HostName }}/{{ Repo.GroupName }}/{{ Repo.ShortName }}/s/{{ Build.Number }}/{{ Stage.Name }}-{{ Build.Number }}-{{ Stage.FinishedTime }}"

	// StepsTransferMarkDemoConfig
	// steps transfer key
	StepsTransferMarkDemoConfig = "demo_config"

	randPasswordCnt  = 8
	randPasswordSeed = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_-!"
)

var (
	// pluginDistTypeSupport
	pluginDistTypeSupport = []string{
		DistTypeGit,
		DistTypeCustom,
	}
)

type (
	// Config file_browser_upload private config
	Config struct {
		Debug             bool
		TimeoutSecond     uint
		StepsTransferPath string
		RootPath          string

		DryRun bool

		FileBrowserBaseConfig FileBrowserBaseConfig

		FileBrowserSendConfig FileBrowserSendConfig
	}

	FileBrowserBaseConfig struct {
		FileBrowserHost              string
		FileBrowserUsername          string
		FileBrowserUserPassword      string
		FileBrowserTimeoutPushSecond uint
		FileBrowserWorkSpace         string
	}

	FileBrowserSendConfig struct {
		FileBrowserRemoteRootPath     string
		FileBrowserDistType           string
		FileBrowserDistGraph          string
		FileBrowserTargetDistRootPath string

		FileBrowserTargetFileGlob    []string
		FileBrowserTargetFileRegular string
		FileBrowserShareLinkEnable   bool
		// FileBrowserShareLinkUnit
		// use
		// [ web_api.ShareUnitDays web_api.ShareUnitHours
		// web_api.ShareUnitMinutes
		// web_api.ShareUnitSeconds ]
		FileBrowserShareLinkUnit               string
		FileBrowserShareLinkExpires            uint
		FileBrowserShareLinkAutoPasswordEnable bool
		FileBrowserShareLinkPassword           string
	}
)
