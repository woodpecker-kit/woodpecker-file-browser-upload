package file_browser_upload_test

import (
	"fmt"
	"github.com/sinlov-go/unittest-kit/env_kit"
	"github.com/sinlov-go/unittest-kit/unittest_file_kit"
	"github.com/woodpecker-kit/woodpecker-file-browser-upload/file_browser_upload"
	"github.com/woodpecker-kit/woodpecker-tools/wd_flag"
	"github.com/woodpecker-kit/woodpecker-tools/wd_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_log"
	"github.com/woodpecker-kit/woodpecker-tools/wd_steps_transfer"
	"github.com/woodpecker-kit/woodpecker-tools/wd_template"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

const (
	keyEnvDebug  = "CI_DEBUG"
	keyEnvCiNum  = "CI_NUMBER"
	keyEnvCiKey  = "CI_KEY"
	keyEnvCiKeys = "CI_KEYS"

	keyLinkSpeedTestUrls = "ENV_LINK_SPEED_TEST_URLS"

	mockVersion = "v1.0.0"
	mockName    = "woodpecker-file-browser-upload"

	mockFileBrowserHost     = "https://file-browser.foo.com"
	mockFileBrowserUserName = "admin"
	mockFileBrowserUserPass = "adminPwd"

	mockFileBrowserRemoteRootPath        = "dist/"
	mockFileBrowserTargetDistRootPath    = ""
	mockFileBrowserTargetFileRegularFail = "*.json"
	mockFileBrowserTargetFileRegularJson = ".*.json"
	mockFileBrowserTargetFileRegularApk  = ".*.apk"
)

var (
	// testBaseFolderPath
	//  test base dir will auto get by package init()
	testBaseFolderPath = ""
	testGoldenKit      *unittest_file_kit.TestGoldenKit

	envTimeoutSecond uint

	// mustSetInCiEnvList
	//  for check set in CI env not empty
	mustSetInCiEnvList = []string{
		wd_flag.EnvKeyCiSystemPlatform,
		wd_flag.EnvKeyCiSystemVersion,
	}
	// mustSetArgsAsEnvList
	mustSetArgsAsEnvList = []string{
		file_browser_upload.EnvFileBrowserHost,
		file_browser_upload.EnvFileBrowserUsername,
		file_browser_upload.EnvFileBrowserUserPassword,
	}

	mockFileGlob = []string{
		"**/*.json",
		"**/*.apk",
	}

	valEnvPluginDebug = false

	valEnvFileBrowserHost              = ""
	valEnvFileBrowserUserName          = ""
	valEnvFileBrowserPassword          = ""
	valEnvFileBrowserTimeoutPushSecond = uint(60)

	valEnvFileGlob = mockFileGlob
)

func init() {
	testBaseFolderPath, _ = getCurrentFolderPath()
	wd_log.SetLogLineDeep(2)

	wd_template.RegisterSettings(wd_template.DefaultHelpers)

	envTimeoutSecond = uint(env_kit.FetchOsEnvInt(wd_flag.EnvKeyPluginTimeoutSecond, 10))

	testGoldenKit = unittest_file_kit.NewTestGoldenKit(testBaseFolderPath)

	valEnvPluginDebug = env_kit.FetchOsEnvBool(wd_flag.EnvKeyPluginDebug, false)
	valEnvFileBrowserHost = env_kit.FetchOsEnvStr(file_browser_upload.EnvFileBrowserHost, "")
	valEnvFileBrowserUserName = env_kit.FetchOsEnvStr(file_browser_upload.EnvFileBrowserUsername, "")
	valEnvFileBrowserPassword = env_kit.FetchOsEnvStr(file_browser_upload.EnvFileBrowserUserPassword, "")
	valEnvFileBrowserTimeoutPushSecond = env_kit.FetchOsEnvUint(file_browser_upload.EnvFileBrowserTimeOutSendSecond, 60)
	valEnvFileGlob = env_kit.FetchOsEnvStringSlice(file_browser_upload.EnvFileBrowserFileGlob)
}

// test case basic tools start
// getCurrentFolderPath
//
//	can get run path this golang dir
func getCurrentFolderPath() (string, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("can not get current file info")
	}
	return filepath.Dir(file), nil
}

// test case basic tools end

func envCheck(t *testing.T) bool {

	if valEnvPluginDebug {
		wd_log.OpenDebug()
	}

	// most CI system will set env CI to true
	envCI := env_kit.FetchOsEnvStr("CI", "")
	if envCI == "" {
		t.Logf("not in CI system, skip envCheck")
		return false
	}
	t.Logf("check env for CI system")
	return env_kit.MustHasEnvSetByArray(t, mustSetInCiEnvList)
}

func envMustArgsCheck(t *testing.T) bool {
	for _, item := range mustSetArgsAsEnvList {
		if os.Getenv(item) == "" {
			t.Logf("plasee set env: %s, than run test\nfull need set env %v", item, mustSetArgsAsEnvList)
			return true
		}
	}
	return false
}

func generateTransferStepsOut(plugin file_browser_upload.FileBrowserPlugin, mark string, data interface{}) error {
	_, err := wd_steps_transfer.Out(plugin.Settings.RootPath, plugin.Settings.StepsTransferPath, plugin.GetWoodPeckerInfo(), mark, data)
	return err
}

func mockPluginSettings() file_browser_upload.Settings {
	// all mock settings can set here
	settings := file_browser_upload.Settings{
		// use env:PLUGIN_DEBUG
		Debug:             valEnvPluginDebug,
		TimeoutSecond:     envTimeoutSecond,
		RootPath:          testGoldenKit.GetTestDataFolderFullPath(),
		StepsTransferPath: wd_steps_transfer.DefaultKitStepsFileName,
	}

	// remove or change this code
	settings.FileBrowserBaseConfig = file_browser_upload.FileBrowserBaseConfig{
		FileBrowserHost:              valEnvFileBrowserHost,
		FileBrowserUsername:          valEnvFileBrowserUserName,
		FileBrowserUserPassword:      valEnvFileBrowserPassword,
		FileBrowserTimeoutPushSecond: valEnvFileBrowserTimeoutPushSecond,
	}
	settings.FileBrowserSendConfig = file_browser_upload.FileBrowserSendConfig{
		FileBrowserRemoteRootPath:     mockFileBrowserRemoteRootPath,
		FileBrowserDistType:           file_browser_upload.DistTypeGit,
		FileBrowserDistGraph:          file_browser_upload.EnvFileBrowserDistGraph,
		FileBrowserTargetDistRootPath: mockFileBrowserTargetDistRootPath,
		FileBrowserTargetFileRegular:  mockFileBrowserTargetFileRegularJson,
	}

	return settings

}

func mockPluginWithSettings(t *testing.T,
	woodpeckerInfo wd_info.WoodpeckerInfo,
	settings file_browser_upload.Settings,
	ciWorkspace string,
) file_browser_upload.FileBrowserPlugin {
	p := file_browser_upload.FileBrowserPlugin{
		Name:    mockName,
		Version: mockVersion,
	}

	if ciWorkspace != "" {
		woodpeckerInfo.BasicInfo.CIWorkspace = ciWorkspace
	}

	// mock woodpecker info
	//t.Log("mockPluginWithStatus")
	p.SetWoodpeckerInfo(woodpeckerInfo)

	p.Settings = settings
	return p
}
