package file_browser_upload

const (
	TimeoutSecondMinimum       = 10
	UploadTimeoutSecondMinimum = 60
	TryConnectTimeoutSecond    = 5
	TryConnectRetries          = 3

	DistTypeGit    = "git"
	DistTypeCustom = "custom"

	// DistGraphTypeGit
	// template is used wd_short_info.WoodpeckerInfoShort
	DistGraphTypeGit = "{{ Repo.Hostname }}/{{ Repo.OwnerName }}/{{ Repo.ShortName }}/s/{{ Build.Number }}/{{ Build.Number }}-{{ Stage.Finished }}"

	distGitGraphDefault     = "{{ Repo.Hostname }}/{{ Repo.OwnerName }}/{{ Repo.ShortName }}/b/{{ Build.Number }}/{{ Commit.Branch }}"
	distGitGraphPullRequest = "{{ Repo.Hostname }}/{{ Repo.OwnerName }}/{{ Repo.ShortName }}/pr/{{ Build.PR }}/{{ Build.Number }}"
	distGitGraphTag         = "{{ Repo.Hostname }}/{{ Repo.OwnerName }}/{{ Repo.ShortName }}/tag/{{ Build.Tag }}/{{ Build.Number }}"

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
	// Settings file_browser_upload private config
	Settings struct {
		Debug             bool
		TimeoutSecond     uint
		StepsTransferPath string
		StepsOutDisable   bool
		RootPath          string

		DryRun bool

		FileBrowserBaseConfig FileBrowserBaseConfig

		FileBrowserSendConfig FileBrowserSendConfig
	}

	FileBrowserBaseConfig struct {
		FileBrowserHost         string
		FileBrowserUsername     string
		FileBrowserUserPassword string

		FileBrowserUrls []string

		FileBrowserTimeoutPushSecond uint
		FileBrowserWorkSpace         string

		usedFileBrowserUrl          string
		usedFileBrowserUsername     string
		usedFileBrowserUserPassword string
	}

	FileBrowserSendConfig struct {
		FileBrowserRemoteRootPath string
		FileBrowserDistType       string
		// FileBrowserDistGraph
		// sample is DistGraphTypeGit, template is used wd_short_info.WoodpeckerInfoShort
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
