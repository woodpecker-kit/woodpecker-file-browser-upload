package file_browser_upload_test

import (
	"encoding/json"
	"github.com/woodpecker-kit/woodpecker-file-browser-upload/file_browser_upload"
	"github.com/woodpecker-kit/woodpecker-tools/wd_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_log"
	"github.com/woodpecker-kit/woodpecker-tools/wd_mock"
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
			if !tc.wantArgFlagNotErr {
				t.Logf("check args error: %v", errPluginRun)
			}
			if (errPluginRun != nil) == tc.wantArgFlagNotErr {
				wd_log.VerboseJsonf(tc.p.Config, "print Config")
				t.Fatalf("Exec() error = %v, wantErr %v", errPluginRun, tc.wantArgFlagNotErr)
				return
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

	// sendSuccessByGlob
	var sendSuccessByGlob file_browser_upload.FileBrowserPlugin
	deepCopyByPlugin(&p, &sendSuccessByGlob)
	sendSuccessByGlob.Config.FileBrowserSendConfig.FileBrowserTargetFileRegular = ""
	sendSuccessByGlob.Config.FileBrowserSendConfig.FileBrowserTargetFileGlob = mockFileGlob

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
			name: "sendSuccessByGlob",
			p:    sendSuccessByGlob,
			//isDryRun: true,
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
		FileBrowserDistGraph:          mockFileBrowserDistGraph,
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
