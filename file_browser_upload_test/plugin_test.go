package file_browser_upload_test

import (
	"encoding/json"
	"github.com/sinlov/filebrowser-client/web_api"
	"github.com/woodpecker-kit/woodpecker-file-browser-upload/file_browser_upload"
	"github.com/woodpecker-kit/woodpecker-tools/wd_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_log"
	"github.com/woodpecker-kit/woodpecker-tools/wd_mock"
	"github.com/woodpecker-kit/woodpecker-tools/wd_short_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_steps_transfer"
	"testing"
)

func TestCheckArgsPlugin(t *testing.T) {
	t.Log("mock FileBrowserPlugin")
	p := mockPluginWithStatus(t, wd_info.BuildStatusSuccess)
	// set default for check args
	p.Config.DryRun = true
	p.Config.FileBrowserBaseConfig.FileBrowserHost = mockFileBrowserHost
	p.Config.FileBrowserBaseConfig.FileBrowserUsername = mockFileBrowserUserName
	p.Config.FileBrowserBaseConfig.FileBrowserUserPassword = mockFileBrowserUserPass

	testDataDistFolderPath, errInitPostFile := initTestDataPostFileDir()
	if errInitPostFile != nil {
		t.Fatal(errInitPostFile)
	}
	p.WoodpeckerInfo.BasicInfo.CIWorkspace = testDataDistFolderPath

	// defaultSettings
	var defaultSettings file_browser_upload.FileBrowserPlugin
	deepCopyByPlugin(&p, &defaultSettings)

	// basicInfoEmptyFileBrowserHost
	var basicInfoEmptyFileBrowserHost file_browser_upload.FileBrowserPlugin
	deepCopyByPlugin(&p, &basicInfoEmptyFileBrowserHost)
	basicInfoEmptyFileBrowserHost.Config.FileBrowserBaseConfig.FileBrowserHost = ""

	// basicInfoEmptyFileBrowserUserName
	var basicInfoEmptyFileBrowserUserName file_browser_upload.FileBrowserPlugin
	deepCopyByPlugin(&p, &basicInfoEmptyFileBrowserUserName)
	basicInfoEmptyFileBrowserUserName.Config.FileBrowserBaseConfig.FileBrowserUsername = ""

	// distTypeNotSupport
	var distTypeNotSupport file_browser_upload.FileBrowserPlugin
	deepCopyByPlugin(&p, &distTypeNotSupport)
	distTypeNotSupport.Config.FileBrowserSendConfig.FileBrowserDistType = "not_support"

	// fileTargetAllEmpty
	var fileTargetAllEmpty file_browser_upload.FileBrowserPlugin
	deepCopyByPlugin(&p, &fileTargetAllEmpty)
	fileTargetAllEmpty.Config.FileBrowserSendConfig.FileBrowserTargetFileRegular = ""
	fileTargetAllEmpty.Config.FileBrowserSendConfig.FileBrowserTargetFileGlob = []string{}

	tests := []struct {
		name              string
		p                 file_browser_upload.FileBrowserPlugin
		isDryRun          bool
		workRoot          string
		wantArgFlagNotErr bool
	}{
		{
			name:              "defaultSettings",
			p:                 defaultSettings,
			wantArgFlagNotErr: true,
		},
		{
			name: "basicInfoEmptyFileBrowserHost",
			p:    basicInfoEmptyFileBrowserHost,
		},
		{
			name: "basicInfoEmptyFileBrowserUserName",
			p:    basicInfoEmptyFileBrowserUserName,
		},
		{
			name: "distTypeNotSupport",
			p:    distTypeNotSupport,
		},
		{
			name: "fileTargetAllEmpty",
			p:    fileTargetAllEmpty,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errPluginRun := tc.p.Exec()
			if tc.wantArgFlagNotErr {
				if errPluginRun != nil {
					wdShotInfo := wd_short_info.ParseWoodpeckerInfo2Short(*tc.p.WoodpeckerInfo)
					wd_log.VerboseJsonf(wdShotInfo, "print WoodpeckerInfoShort")
					wd_log.VerboseJsonf(tc.p.Config, "print Config")
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

func TestPlugin(t *testing.T) {
	t.Log("do FileBrowserPlugin")
	if envCheck(t) {
		return
	}
	if envMustArgsCheck(t) {
		return
	}
	t.Log("mock FileBrowserPlugin")
	p := mockPluginWithStatus(t, wd_info.BuildStatusSuccess)

	testDataDistFolderPath, errInitPostFile := initTestDataPostFileDir()
	if errInitPostFile != nil {
		t.Fatal(errInitPostFile)
	}
	p.WoodpeckerInfo.BasicInfo.CIWorkspace = testDataDistFolderPath

	//wd_log.VerboseJsonf(p, "print file_browser_upload info")

	t.Log("mock file_browser_upload config")

	// sendSuccessByRegularJson
	var sendSuccessByRegularJson file_browser_upload.FileBrowserPlugin
	deepCopyByPlugin(&p, &sendSuccessByRegularJson)

	// sendCustomByGlob
	var sendCustomByGlob file_browser_upload.FileBrowserPlugin
	deepCopyByPlugin(&p, &sendCustomByGlob)
	sendCustomByGlob.Config.FileBrowserSendConfig.FileBrowserDistType = file_browser_upload.DistTypeCustom
	sendCustomByGlob.Config.FileBrowserSendConfig.FileBrowserDistGraph = file_browser_upload.DistGraphTypeGit
	sendCustomByGlob.Config.FileBrowserSendConfig.FileBrowserTargetFileRegular = ""
	sendCustomByGlob.Config.FileBrowserSendConfig.FileBrowserTargetFileGlob = mockFileGlob

	// sendSuccessByGlob
	var sendSuccessByGlob file_browser_upload.FileBrowserPlugin
	deepCopyByPlugin(&p, &sendSuccessByGlob)
	sendSuccessByGlob.Config.FileBrowserSendConfig.FileBrowserTargetFileRegular = ""
	sendSuccessByGlob.Config.FileBrowserSendConfig.FileBrowserTargetFileGlob = mockFileGlob

	// sendSuccessPrGlob
	var sendSuccessPrGlob file_browser_upload.FileBrowserPlugin
	deepCopyByPlugin(&p, &sendSuccessPrGlob)
	sendSuccessPrGlob.Config.FileBrowserSendConfig.FileBrowserTargetFileRegular = ""
	sendSuccessPrGlob.Config.FileBrowserSendConfig.FileBrowserTargetFileGlob = mockFileGlob
	sendSuccessPrGlob.WoodpeckerInfo = wd_mock.NewWoodpeckerInfo(
		wd_mock.WithFastMockPullRequest("1", "new pr", "feature-support", "main", "main"),
	)
	sendSuccessPrGlob.WoodpeckerInfo.BasicInfo.CIWorkspace = testGoldenKit.GetTestDataFolderFullPath()

	// sendSuccessTagGlob
	var sendSuccessTagGlob file_browser_upload.FileBrowserPlugin
	deepCopyByPlugin(&p, &sendSuccessTagGlob)
	sendSuccessTagGlob.Config.FileBrowserSendConfig.FileBrowserTargetFileRegular = ""
	sendSuccessTagGlob.Config.FileBrowserSendConfig.FileBrowserTargetFileGlob = mockFileGlob
	sendSuccessTagGlob.WoodpeckerInfo = wd_mock.NewWoodpeckerInfo(
		wd_mock.WithFastMockTag("v1.0.0", "new tag"),
	)
	sendSuccessTagGlob.WoodpeckerInfo.BasicInfo.CIWorkspace = testGoldenKit.GetTestDataFolderFullPath()

	// sendTagGlobWithShare
	var sendTagGlobWithShare file_browser_upload.FileBrowserPlugin
	deepCopyByPlugin(&p, &sendTagGlobWithShare)
	sendTagGlobWithShare.Config.FileBrowserSendConfig.FileBrowserTargetFileRegular = ""
	sendTagGlobWithShare.Config.FileBrowserSendConfig.FileBrowserTargetFileGlob = mockFileGlob
	sendTagGlobWithShare.WoodpeckerInfo = wd_mock.NewWoodpeckerInfo(
		wd_mock.WithFastMockTag("v1.0.0", "new tag"),
	)
	sendTagGlobWithShare.WoodpeckerInfo.BasicInfo.CIWorkspace = testGoldenKit.GetTestDataFolderFullPath()
	sendTagGlobWithShare.Config.FileBrowserSendConfig.FileBrowserShareLinkEnable = true
	sendTagGlobWithShare.Config.FileBrowserSendConfig.FileBrowserShareLinkExpires = 12
	sendTagGlobWithShare.Config.FileBrowserSendConfig.FileBrowserShareLinkUnit = web_api.ShareUnitHours
	sendTagGlobWithShare.Config.FileBrowserSendConfig.FileBrowserShareLinkAutoPasswordEnable = true

	tests := []struct {
		name            string
		p               file_browser_upload.FileBrowserPlugin
		isDryRun        bool
		workRoot        string
		ossTransferKey  string
		ossTransferData interface{}
		wantErr         bool
	}{
		{
			name:     "sendSuccessByRegularJson",
			p:        sendSuccessByRegularJson,
			isDryRun: true,
		},
		{
			name:     "sendCustomByGlob",
			p:        sendCustomByGlob,
			isDryRun: true,
		},
		{
			name:     "sendSuccessByGlob",
			p:        sendSuccessByGlob,
			isDryRun: true,
		},
		{
			name:     "sendSuccessPrGlob",
			p:        sendSuccessPrGlob,
			isDryRun: true,
		},
		{
			name:     "sendSuccessTagGlob",
			p:        sendSuccessTagGlob,
			isDryRun: true,
		},
		{
			name:     "sendTagGlobWithShare",
			p:        sendTagGlobWithShare,
			isDryRun: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.p.Config.DryRun = tc.isDryRun
			if tc.workRoot != "" {
				tc.p.Config.RootPath = tc.workRoot
				errGenTransferData := generateTransferStepsOut(
					tc.p,
					tc.ossTransferKey,
					tc.ossTransferData,
				)
				if errGenTransferData != nil {
					t.Fatal(errGenTransferData)
				}
			}
			err := tc.p.Exec()
			if (err != nil) != tc.wantErr {
				t.Errorf("FeishuPlugin.Exec() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}

func mockPluginWithStatus(t *testing.T, status string) file_browser_upload.FileBrowserPlugin {
	p := file_browser_upload.FileBrowserPlugin{
		Name:    mockName,
		Version: mockVersion,
	}
	// use env:PLUGIN_DEBUG
	p.Config.Debug = valEnvPluginDebug
	p.Config.TimeoutSecond = envTimeoutSecond
	p.Config.RootPath = testGoldenKit.GetTestDataFolderFullPath()
	p.Config.StepsTransferPath = wd_steps_transfer.DefaultKitStepsFileName

	// mock woodpecker info
	//t.Log("mockPluginWithStatus")
	woodpeckerInfo := wd_mock.NewWoodpeckerInfo(
		wd_mock.WithCurrentPipelineStatus(status),
	)
	p.WoodpeckerInfo = woodpeckerInfo
	p.WoodpeckerInfo.BasicInfo.CIWorkspace = testGoldenKit.GetTestDataFolderFullPath()

	p.Config.FileBrowserBaseConfig = file_browser_upload.FileBrowserBaseConfig{
		FileBrowserHost:              valEnvFileBrowserHost,
		FileBrowserUsername:          valEnvFileBrowserUserName,
		FileBrowserUserPassword:      valEnvFileBrowserPassword,
		FileBrowserTimeoutPushSecond: valEnvFileBrowserTimeoutPushSecond,
	}

	p.Config.FileBrowserSendConfig = file_browser_upload.FileBrowserSendConfig{
		FileBrowserRemoteRootPath:     mockFileBrowserRemoteRootPath,
		FileBrowserDistType:           file_browser_upload.DistTypeGit,
		FileBrowserDistGraph:          file_browser_upload.EnvFileBrowserDistGraph,
		FileBrowserTargetDistRootPath: mockFileBrowserTargetDistRootPath,
		FileBrowserTargetFileRegular:  mockFileBrowserTargetFileRegularJson,
	}

	return p
}

func deepCopyByPlugin(src, dst *file_browser_upload.FileBrowserPlugin) {
	if tmp, err := json.Marshal(&src); err != nil {
		return
	} else {
		err = json.Unmarshal(tmp, dst)
		return
	}
}

func generateTransferStepsOut(plugin file_browser_upload.FileBrowserPlugin, mark string, data interface{}) error {
	_, err := wd_steps_transfer.Out(plugin.Config.RootPath, plugin.Config.StepsTransferPath, *plugin.WoodpeckerInfo, mark, data)
	return err
}
