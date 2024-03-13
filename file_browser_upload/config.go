package file_browser_upload

const (
	DistTypeGit    = "git"
	DistTypeCustom = "custom"

	// DistGraphTypeGit
	// template is used wd_info_shot.WoodpeckerInfoShort
	DistGraphTypeGit = "{{ Repo.HostName }}/{{ Repo.OwnerName }}/{{ Repo.ShortName }}/s/{{ Build.Number }}/{{ Build.Number }}-{{ Stage.Finished }}"

	distGitGraphDefault     = "{{ Repo.HostName }}/{{ Repo.OwnerName }}/{{ Repo.ShortName }}/b/{{ Build.Number }}/{{ Commit.Branch }}"
	distGitGraphPullRequest = "{{ Repo.HostName }}/{{ Repo.OwnerName }}/{{ Repo.ShortName }}/pr/{{ Build.PR }}/{{ Build.Number }}"
	distGitGraphTag         = "{{ Repo.HostName }}/{{ Repo.OwnerName }}/{{ Repo.ShortName }}/tag/{{ Build.Tag }}/{{ Build.Number }}"

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
		StepsOutDisable   bool
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
		FileBrowserRemoteRootPath string
		FileBrowserDistType       string
		// FileBrowserDistGraph
		// sample is DistGraphTypeGit, template is used wd_info_shot.WoodpeckerInfoShort
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
