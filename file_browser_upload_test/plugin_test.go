package file_browser_upload_test

import (
	"github.com/sinlov/filebrowser-client/web_api"
	"github.com/woodpecker-kit/woodpecker-file-browser-upload/file_browser_upload"
	"github.com/woodpecker-kit/woodpecker-tools/wd_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_log"
	"github.com/woodpecker-kit/woodpecker-tools/wd_mock"
	"github.com/woodpecker-kit/woodpecker-tools/wd_short_info"
	"testing"
)

func TestCheckArgsPlugin(t *testing.T) {
	t.Log("mock FileBrowserPlugin")
	// set default for check args

	testDataDistFolderPath, errInitPostFile := initTestDataPostFileDir()
	if errInitPostFile != nil {
		t.Fatal(errInitPostFile)
	}

	// defaultSettings
	defaultSettingsWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.WithCurrentPipelineStatus(wd_info.BuildStatusSuccess),
	)
	defaultSettings := mockForTestArgsPluginSettings()

	// basicInfoEmptyFileBrowserHost
	basicInfoEmptyFileBrowserHostWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.WithCurrentPipelineStatus(wd_info.BuildStatusSuccess),
	)
	basicInfoEmptyFileBrowserHostSettings := mockForTestArgsPluginSettings()
	basicInfoEmptyFileBrowserHostSettings.FileBrowserBaseConfig.FileBrowserHost = ""

	// basicInfoEmptyFileBrowserUserName
	basicInfoEmptyFileBrowserUserNameWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.WithCurrentPipelineStatus(wd_info.BuildStatusSuccess),
	)
	basicInfoEmptyFileBrowserUserNameSettings := mockForTestArgsPluginSettings()
	basicInfoEmptyFileBrowserUserNameSettings.FileBrowserBaseConfig.FileBrowserUsername = ""

	// distTypeNotSupport
	distTypeNotSupportWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.WithCurrentPipelineStatus(wd_info.BuildStatusSuccess),
	)
	distTypeNotSupportSettings := mockForTestArgsPluginSettings()
	distTypeNotSupportSettings.FileBrowserSendConfig.FileBrowserDistType = "not_support"

	// fileTargetAllEmpty
	fileTargetAllEmptyWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.WithCurrentPipelineStatus(wd_info.BuildStatusSuccess),
	)
	fileTargetAllEmptySettings := mockForTestArgsPluginSettings()
	fileTargetAllEmptySettings.FileBrowserSendConfig.FileBrowserTargetFileRegular = ""
	fileTargetAllEmptySettings.FileBrowserSendConfig.FileBrowserTargetFileGlob = []string{}

	tests := []struct {
		name              string
		woodpeckerInfo    wd_info.WoodpeckerInfo
		settings          file_browser_upload.Settings
		ciWorkspace       string
		isDryRun          bool
		workRoot          string
		wantArgFlagNotErr bool
	}{
		{
			name:              "defaultSettings",
			woodpeckerInfo:    defaultSettingsWoodpeckerInfo,
			settings:          defaultSettings,
			ciWorkspace:       testDataDistFolderPath,
			wantArgFlagNotErr: true,
		},
		{
			name:           "basicInfoEmptyFileBrowserHost",
			woodpeckerInfo: basicInfoEmptyFileBrowserHostWoodpeckerInfo,
			settings:       basicInfoEmptyFileBrowserHostSettings,
			ciWorkspace:    testDataDistFolderPath,
		},
		{
			name:           "basicInfoEmptyFileBrowserUserName",
			woodpeckerInfo: basicInfoEmptyFileBrowserUserNameWoodpeckerInfo,
			settings:       basicInfoEmptyFileBrowserUserNameSettings,
			ciWorkspace:    testDataDistFolderPath,
		},
		{
			name:           "distTypeNotSupport",
			woodpeckerInfo: distTypeNotSupportWoodpeckerInfo,
			settings:       distTypeNotSupportSettings,
			ciWorkspace:    testDataDistFolderPath,
		},
		{
			name:           "fileTargetAllEmpty",
			woodpeckerInfo: fileTargetAllEmptyWoodpeckerInfo,
			settings:       fileTargetAllEmptySettings,
			ciWorkspace:    testDataDistFolderPath,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := mockPluginWithSettings(t, tc.woodpeckerInfo, tc.settings, tc.ciWorkspace)
			p.OnlyArgsCheck()
			errPluginRun := p.Exec()
			if tc.wantArgFlagNotErr {
				if errPluginRun != nil {
					wdShotInfo := wd_short_info.ParseWoodpeckerInfo2Short(p.GetWoodPeckerInfo())
					wd_log.VerboseJsonf(wdShotInfo, "print WoodpeckerInfoShort")
					wd_log.VerboseJsonf(p.Settings, "print Settings")
					t.Fatalf("wantArgFlagNotErr %v\np.Exec() error:\n%v", tc.wantArgFlagNotErr, errPluginRun)
					return
				}
			} else {
				if errPluginRun == nil {
					t.Fatalf("test case [ %s ], wantArgFlagNotErr %v, but p.Exec() not error", tc.name, tc.wantArgFlagNotErr)
				}
				t.Logf("check args error: %v", errPluginRun)
			}
		})
	}
}

func mockForTestArgsPluginSettings() file_browser_upload.Settings {
	settings := mockPluginSettings()

	settings.FileBrowserBaseConfig.FileBrowserHost = mockFileBrowserHost
	settings.FileBrowserBaseConfig.FileBrowserUsername = mockFileBrowserUserName
	settings.FileBrowserBaseConfig.FileBrowserUserPassword = mockFileBrowserUserPass

	return settings

}

func TestPlugin(t *testing.T) {
	t.Log("do FileBrowserPlugin")
	if envCheck(t) {
		return
	}
	if envMustArgsCheck(t) {
		return
	}
	t.Log("mock FileBrowserPlugin")

	testDataDistFolderPath, errInitPostFile := initTestDataPostFileDir()
	if errInitPostFile != nil {
		t.Fatal(errInitPostFile)
	}

	t.Log("mock file_browser_upload config")

	// sendSuccessByRegularJson
	sendSuccessByRegularJsonWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.WithCurrentPipelineStatus(wd_info.BuildStatusSuccess),
	)
	sendSuccessByRegularJsonSettings := mockPluginSettings()

	// sendCustomByGlob
	sendCustomByGlobWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.WithCurrentPipelineStatus(wd_info.BuildStatusSuccess),
	)
	sendCustomByGlobSettings := mockPluginSettings()
	sendCustomByGlobSettings.FileBrowserSendConfig.FileBrowserDistType = file_browser_upload.DistTypeCustom
	sendCustomByGlobSettings.FileBrowserSendConfig.FileBrowserDistGraph = file_browser_upload.DistGraphTypeGit
	sendCustomByGlobSettings.FileBrowserSendConfig.FileBrowserTargetFileRegular = ""
	sendCustomByGlobSettings.FileBrowserSendConfig.FileBrowserTargetFileGlob = mockFileGlob

	// sendSuccessByGlob
	sendSuccessByGlobalWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.WithCurrentPipelineStatus(wd_info.BuildStatusSuccess),
	)
	sendSuccessByGlobalSettings := mockPluginSettings()
	sendSuccessByGlobalSettings.FileBrowserSendConfig.FileBrowserTargetFileRegular = ""
	sendSuccessByGlobalSettings.FileBrowserSendConfig.FileBrowserTargetFileGlob = mockFileGlob

	// sendSuccessPrGlob
	sendSuccessPrGlobalWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.WithFastMockPullRequest("1", "new pr", "feature-support", "main", "main"),
	)
	sendSuccessPrGlobalSettings := mockPluginSettings()
	sendSuccessPrGlobalSettings.FileBrowserSendConfig.FileBrowserTargetFileRegular = ""
	sendSuccessPrGlobalSettings.FileBrowserSendConfig.FileBrowserTargetFileGlob = mockFileGlob

	// sendSuccessTagGlob
	sendSuccessTagGlobalWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.WithFastMockTag("v1.0.0", "new tag"),
	)
	sendSuccessTagGlobalSettings := mockPluginSettings()
	sendSuccessTagGlobalSettings.FileBrowserSendConfig.FileBrowserTargetFileRegular = ""
	sendSuccessTagGlobalSettings.FileBrowserSendConfig.FileBrowserTargetFileGlob = mockFileGlob

	// sendTagGlobWithShare
	sendTagGlobWithShareWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.WithFastMockTag("v1.0.0", "new tag"),
	)
	sendTagGlobWithShareSettings := mockPluginSettings()
	sendTagGlobWithShareSettings.FileBrowserSendConfig.FileBrowserTargetFileRegular = ""
	sendTagGlobWithShareSettings.FileBrowserSendConfig.FileBrowserTargetFileGlob = mockFileGlob
	sendTagGlobWithShareSettings.FileBrowserSendConfig.FileBrowserShareLinkEnable = true
	sendTagGlobWithShareSettings.FileBrowserSendConfig.FileBrowserShareLinkExpires = 12
	sendTagGlobWithShareSettings.FileBrowserSendConfig.FileBrowserShareLinkUnit = web_api.ShareUnitHours
	sendTagGlobWithShareSettings.FileBrowserSendConfig.FileBrowserShareLinkAutoPasswordEnable = true

	tests := []struct {
		name           string
		woodpeckerInfo wd_info.WoodpeckerInfo
		settings       file_browser_upload.Settings
		ciWorkspace    string
		isDryRun       bool
		workRoot       string

		ossTransferKey  string
		ossTransferData interface{}
		wantErr         bool
	}{
		{
			name:           "sendSuccessByRegularJson",
			woodpeckerInfo: sendSuccessByRegularJsonWoodpeckerInfo,
			settings:       sendSuccessByRegularJsonSettings,
			ciWorkspace:    testDataDistFolderPath,
			isDryRun:       true,
		},
		{
			name:           "sendCustomByGlob",
			woodpeckerInfo: sendCustomByGlobWoodpeckerInfo,
			settings:       sendCustomByGlobSettings,
			ciWorkspace:    testDataDistFolderPath,
			isDryRun:       true,
		},
		{
			name:           "sendSuccessByGlob",
			woodpeckerInfo: sendSuccessByGlobalWoodpeckerInfo,
			settings:       sendSuccessByGlobalSettings,
			ciWorkspace:    testDataDistFolderPath,
			isDryRun:       true,
		},
		{
			name:           "sendSuccessPrGlob",
			woodpeckerInfo: sendSuccessPrGlobalWoodpeckerInfo,
			settings:       sendSuccessPrGlobalSettings,
			ciWorkspace:    testDataDistFolderPath,
			isDryRun:       true,
		},
		{
			name:           "sendSuccessTagGlob",
			woodpeckerInfo: sendSuccessTagGlobalWoodpeckerInfo,
			settings:       sendSuccessTagGlobalSettings,
			ciWorkspace:    testDataDistFolderPath,
			isDryRun:       true,
		},
		{
			name:           "sendTagGlobWithShare",
			woodpeckerInfo: sendTagGlobWithShareWoodpeckerInfo,
			settings:       sendTagGlobWithShareSettings,
			ciWorkspace:    testDataDistFolderPath,
			isDryRun:       false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := mockPluginWithSettings(t, tc.woodpeckerInfo, tc.settings, tc.ciWorkspace)
			p.Settings.DryRun = tc.isDryRun
			if tc.workRoot != "" {
				p.Settings.RootPath = tc.workRoot
				errGenTransferData := generateTransferStepsOut(
					p,
					tc.ossTransferKey,
					tc.ossTransferData,
				)
				if errGenTransferData != nil {
					t.Fatal(errGenTransferData)
				}
			}

			err := p.Exec()
			if (err != nil) != tc.wantErr {
				t.Errorf("FeishuPlugin.Exec() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}
