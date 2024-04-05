package file_browser_upload

import (
	"fmt"
	"github.com/sinlov-go/go-common-lib/pkg/folder"
	"github.com/sinlov-go/go-common-lib/pkg/string_tools"
	"github.com/sinlov-go/go-common-lib/pkg/struct_kit"
	"github.com/sinlov/filebrowser-client/file_browser_client"
	tools "github.com/sinlov/filebrowser-client/tools/str_tools"
	"github.com/sinlov/filebrowser-client/web_api"
	"github.com/woodpecker-kit/woodpecker-tools/wd_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_log"
	"github.com/woodpecker-kit/woodpecker-tools/wd_short_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_steps_transfer"
	"github.com/woodpecker-kit/woodpecker-tools/wd_template"
	"github.com/woodpecker-kit/woodpecker-transfer-data/wd_share_file_browser_upload"
	"math/rand"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (p *FileBrowserPlugin) ShortInfo() wd_short_info.WoodpeckerInfoShort {
	if p.wdShortInfo == nil {
		info2Short := wd_short_info.ParseWoodpeckerInfo2Short(*p.woodpeckerInfo)
		p.wdShortInfo = &info2Short
	}
	return *p.wdShortInfo
}

// SetWoodpeckerInfo
// also change ShortInfo() return
func (p *FileBrowserPlugin) SetWoodpeckerInfo(info wd_info.WoodpeckerInfo) {
	var newInfo wd_info.WoodpeckerInfo
	_ = struct_kit.DeepCopyByGob(&info, &newInfo)
	p.woodpeckerInfo = &newInfo
	info2Short := wd_short_info.ParseWoodpeckerInfo2Short(newInfo)
	p.wdShortInfo = &info2Short
}

func (p *FileBrowserPlugin) GetWoodPeckerInfo() wd_info.WoodpeckerInfo {
	return *p.woodpeckerInfo
}

func (p *FileBrowserPlugin) OnlyArgsCheck() {
	p.onlyArgsCheck = true
}

func (p *FileBrowserPlugin) Exec() error {
	errLoadStepsTransfer := p.loadStepsTransfer()
	if errLoadStepsTransfer != nil {
		return errLoadStepsTransfer
	}

	errCheckArgs := p.checkArgs()
	if errCheckArgs != nil {
		return fmt.Errorf("check args err: %v", errCheckArgs)
	}

	if p.onlyArgsCheck {
		wd_log.Info("only check args, skip do doBiz")
		return nil
	}

	err := p.doBiz()
	if err != nil {
		return err
	}
	errSaveStepsTransfer := p.saveStepsTransfer()
	if errSaveStepsTransfer != nil {
		return errSaveStepsTransfer
	}

	return nil
}

func (p *FileBrowserPlugin) loadStepsTransfer() error {
	return nil
}

func (p *FileBrowserPlugin) checkArgs() error {
	errCheck := argCheckInArr("args file-browser-dist-type", p.Settings.FileBrowserSendConfig.FileBrowserDistType, pluginDistTypeSupport)
	if errCheck != nil {
		return errCheck
	}

	if len(p.Settings.FileBrowserSendConfig.FileBrowserTargetFileGlob) == 0 && p.Settings.FileBrowserSendConfig.FileBrowserTargetFileRegular == "" {
		return fmt.Errorf("args file-browser-file-glob and file-browser-file-regular not be empty")
	}
	if p.Settings.FileBrowserBaseConfig.FileBrowserHost == "" {
		return fmt.Errorf("args file-browser-host not be empty")
	}

	if p.Settings.FileBrowserBaseConfig.FileBrowserUsername == "" {
		return fmt.Errorf("args file-browser-username not be empty")
	}

	if p.Settings.FileBrowserSendConfig.FileBrowserRemoteRootPath == "" {
		return fmt.Errorf("args file-browser-remote-root-path not be empty")
	}

	// check default FileBrowserTimeoutPushSecond
	if p.Settings.FileBrowserBaseConfig.FileBrowserTimeoutPushSecond < 60 {
		p.Settings.FileBrowserBaseConfig.FileBrowserTimeoutPushSecond = 60
	}

	// check default p.Settings.FileBrowserBaseConfig.FileBrowserWorkSpace
	if p.Settings.FileBrowserBaseConfig.FileBrowserWorkSpace == "" {
		p.Settings.FileBrowserBaseConfig.FileBrowserWorkSpace = p.woodpeckerInfo.BasicInfo.CIWorkspace
	}

	return nil
}

func argCheckInArr(mark string, target string, checkArr []string) error {
	if !(string_tools.StringInArr(target, checkArr)) {
		return fmt.Errorf("not support %s now [ %s ], must in %v", mark, target, checkArr)
	}
	return nil
}

// doBiz
//
//	replace this code with your file_browser_upload implementation
func (p *FileBrowserPlugin) doBiz() error {

	p.shareFileBrowserUpload = &wd_share_file_browser_upload.WdShareFileBrowserUpload{
		IsSendSuccess: false,
	}

	fileBrowserClient, errNew := file_browser_client.NewClient(
		p.Settings.FileBrowserBaseConfig.FileBrowserUsername,
		p.Settings.FileBrowserBaseConfig.FileBrowserUserPassword,
		p.Settings.FileBrowserBaseConfig.FileBrowserHost,
		p.Settings.TimeoutSecond,
		p.Settings.FileBrowserBaseConfig.FileBrowserTimeoutPushSecond,
	)

	if errNew != nil {
		return fmt.Errorf("new fileBrowser client err: %v", errNew)
	}
	fileBrowserClient.Debug(p.Settings.Debug)

	var remoteRealRootPath = strings.TrimRight(p.Settings.FileBrowserSendConfig.FileBrowserRemoteRootPath, "/")

	switch p.Settings.FileBrowserSendConfig.FileBrowserDistType {
	default:
		return fmt.Errorf("send dist type not support %s", p.Settings.FileBrowserSendConfig.FileBrowserDistType)
	case DistTypeGit:
		commitShortSha := string([]rune(p.woodpeckerInfo.CurrentInfo.CurrentCommitInfo.CiCommitSha))[:8]
		if p.woodpeckerInfo.CurrentInfo.CurrentCommitInfo.CiCommitTag != "" {

			tagPath, errPathTag := wd_template.RenderTrim(distGitGraphTag, p.ShortInfo())
			if errPathTag != nil {
				return fmt.Errorf("render as %s \nerr: %v", distGitGraphTag, errPathTag)
			}
			remoteRealRootPath = path.Join(remoteRealRootPath, tagPath, commitShortSha)
		} else if p.woodpeckerInfo.CurrentInfo.CurrentCommitInfo.CiCommitPullRequest != "" {
			prPath, errPathPr := wd_template.RenderTrim(distGitGraphPullRequest, p.ShortInfo())
			if errPathPr != nil {
				return fmt.Errorf("render as %s \nerr: %v", distGitGraphPullRequest, errPathPr)
			}
			remoteRealRootPath = path.Join(remoteRealRootPath, prPath, commitShortSha)
		} else {
			defaultPath, errPathDefault := wd_template.RenderTrim(distGitGraphDefault, p.ShortInfo())
			if errPathDefault != nil {
				return fmt.Errorf("render as %s \nerr: %v", distGitGraphDefault, errPathDefault)
			}

			remoteRealRootPath = path.Join(remoteRealRootPath, defaultPath, commitShortSha)
		}
	case DistTypeCustom:
		renderPath, err := wd_template.RenderTrim(p.Settings.FileBrowserSendConfig.FileBrowserDistGraph, p.ShortInfo())
		if err != nil {
			return fmt.Errorf("setting file-browser-dist-graph as %s \nerr: %v", p.Settings.FileBrowserSendConfig.FileBrowserDistGraph, err)
		}
		remoteRealRootPath = path.Join(remoteRealRootPath, renderPath)
	}

	wd_log.Debugf("final remoteRealRootPath: %s", remoteRealRootPath)

	targetRootPath := filepath.Join(p.Settings.FileBrowserBaseConfig.FileBrowserWorkSpace, p.Settings.FileBrowserSendConfig.FileBrowserTargetDistRootPath)
	wd_log.Debugf("workOnSend fileBrowserDistType: %s", p.Settings.FileBrowserSendConfig.FileBrowserDistType)
	if p.Settings.FileBrowserSendConfig.FileBrowserDistType == DistTypeCustom {
		wd_log.Verbosef("workOnSend fileBrowserDistGraph: %s", p.Settings.FileBrowserSendConfig.FileBrowserDistGraph)
	}
	wd_log.Debugf("workOnSend remoteRealRootPath: %s", remoteRealRootPath)
	wd_log.Debugf("workOnSend targetRootPath: %s", targetRootPath)
	wd_log.Debugf("workOnSend fileBrowserWorkSpace: %s", p.Settings.FileBrowserBaseConfig.FileBrowserWorkSpace)
	wd_log.Debugf("workOnSend targetDistRootPath: %s", p.Settings.FileBrowserSendConfig.FileBrowserTargetDistRootPath)
	wd_log.Debugf("workOnSend target File Regular: %v", p.Settings.FileBrowserSendConfig.FileBrowserTargetFileRegular)
	wd_log.Debugf("workOnSend target File Glob: %v", p.Settings.FileBrowserSendConfig.FileBrowserTargetFileGlob)

	if !(folder.PathExistsFast(targetRootPath)) {
		return fmt.Errorf("file browser want send file local path not exists at: %s", targetRootPath)
	}

	var fileSendPathList []string
	if folder.PathIsFile(targetRootPath) {
		wd_log.Debugf("target path is file just send one local file: %s", targetRootPath)
		fileSendPathList = append(fileSendPathList, targetRootPath)
	} else {
		if p.Settings.FileBrowserSendConfig.FileBrowserTargetFileGlob != nil && len(p.Settings.FileBrowserSendConfig.FileBrowserTargetFileGlob) > 0 {
			wd_log.Debugf("target file want find by File Glob: %v", p.Settings.FileBrowserSendConfig.FileBrowserTargetFileGlob)
			for _, glob := range p.Settings.FileBrowserSendConfig.FileBrowserTargetFileGlob {
				walkByGlob, errWalkAllByGlob := folder.WalkAllByGlob(targetRootPath, glob, true)
				if errWalkAllByGlob != nil {
					return fmt.Errorf("file browser want send file local path with glob %s be err: %v", targetRootPath, errWalkAllByGlob)
				}
				wd_log.Debugf("target path find by File Glob [ %s ] files:\n%s", glob, strings.Join(walkByGlob, "\n"))
				fileSendPathList = append(fileSendPathList, walkByGlob...)
			}
		}
		if p.Settings.FileBrowserSendConfig.FileBrowserTargetFileRegular != "" {
			wd_log.Debugf("target file want find by File Regular: %v", p.Settings.FileBrowserSendConfig.FileBrowserTargetFileRegular)

			matchPath, err := folder.WalkAllByMatchPath(targetRootPath, p.Settings.FileBrowserSendConfig.FileBrowserTargetFileRegular, true)
			if err != nil {
				return fmt.Errorf("file browser want send file local path with file regular %s be err: %v", targetRootPath, err)
			}
			wd_log.Debugf("target path find by File Regular [ %s ] files:\n%s", p.Settings.FileBrowserSendConfig.FileBrowserTargetFileRegular, strings.Join(matchPath, "\n"))
			fileSendPathList = append(fileSendPathList, matchPath...)
		}
	}

	if len(fileSendPathList) == 0 {
		return fmt.Errorf("file browser want send file local path not find any file at path: %s\nglob: %s\nregular: %s",
			targetRootPath,
			p.Settings.FileBrowserSendConfig.FileBrowserTargetFileGlob,
			p.Settings.FileBrowserSendConfig.FileBrowserTargetFileRegular,
		)
	}

	wd_log.Debugf("now send path len %d", len(fileSendPathList))
	fileSendPathList = tools.StrArrRemoveDuplicates(fileSendPathList)
	if p.Settings.Debug {
		wd_log.Debugf("debug: send path remove duplicates len %d", len(fileSendPathList))
	}

	if p.Settings.DryRun {
		wd_log.Infof("dry run mode not send file to file browser, more info to open debug")
		return nil
	}

	err := fileBrowserClient.Login()
	if err != nil {
		return err
	}

	if len(fileSendPathList) == 1 {
		localFileAbsPath := fileSendPathList[0]
		remotePath := fetchRemotePathByLocalRoot(localFileAbsPath, targetRootPath, remoteRealRootPath)
		var resourcePostOne = file_browser_client.ResourcePostFile{
			LocalFilePath:  localFileAbsPath,
			RemoteFilePath: remotePath,
		}
		errSendOneFile := fileBrowserClient.ResourcesPostFile(resourcePostOne, p.Settings.Debug)
		if errSendOneFile != nil {
			return errSendOneFile
		}
		if p.Settings.FileBrowserSendConfig.FileBrowserShareLinkEnable {
			errSendFileShare := shareBySendConfig(fileBrowserClient, p, remotePath, false)
			if errSendFileShare != nil {
				return errSendFileShare
			}
		}

	} else {
		for _, item := range fileSendPathList {
			var resourcePost = file_browser_client.ResourcePostFile{
				LocalFilePath:  item,
				RemoteFilePath: fetchRemotePathByLocalRoot(item, targetRootPath, remoteRealRootPath),
			}
			errSendOneFile := fileBrowserClient.ResourcesPostFile(resourcePost, p.Settings.Debug)
			if errSendOneFile != nil {
				return errSendOneFile
			}
		}
		if p.Settings.FileBrowserSendConfig.FileBrowserShareLinkEnable {
			errSendFileShare := shareBySendConfig(fileBrowserClient, p, remoteRealRootPath, true)
			if errSendFileShare != nil {
				return errSendFileShare
			}
		}
	}

	return nil
}

func fetchRemotePathByLocalRoot(localAbsPath, localRootPath, remoteRootPath string) string {
	remotePath := strings.Replace(localAbsPath, localRootPath, "", -1)
	remotePath = folder.Path2WebPath(remotePath)
	return fmt.Sprintf("%s/%s", remoteRootPath, remotePath)
}

func shareBySendConfig(client file_browser_client.FileBrowserClient, p *FileBrowserPlugin, remotePath string, isDir bool) error {
	expires := strconv.Itoa(int(p.Settings.FileBrowserSendConfig.FileBrowserShareLinkExpires))
	passWord := p.Settings.FileBrowserSendConfig.FileBrowserShareLinkPassword
	if p.Settings.FileBrowserSendConfig.FileBrowserShareLinkAutoPasswordEnable {
		passWord = genPwd(randPasswordCnt)
	}
	if isDir {
		remotePath = fmt.Sprintf("%s/", remotePath)
	}
	shareResource := file_browser_client.ShareResource{
		RemotePath: remotePath,
		ShareConfig: web_api.ShareConfig{
			Password: passWord,
			Expires:  expires,
			Unit:     p.Settings.FileBrowserSendConfig.FileBrowserShareLinkUnit,
		},
	}
	sharePost, errSendShareFile := client.SharePost(shareResource)
	if errSendShareFile != nil {
		return errSendShareFile
	}
	wd_log.Infof("=> share page: %s", sharePost.DownloadPage)
	var shareFileBrowserUpload wd_share_file_browser_upload.WdShareFileBrowserUpload
	if passWord == "" {
		shareFileBrowserUpload = wd_share_file_browser_upload.WdShareFileBrowserUpload{
			IsSendSuccess:       true,
			HostUrl:             p.Settings.FileBrowserBaseConfig.FileBrowserHost,
			FileBrowserUserName: p.Settings.FileBrowserBaseConfig.FileBrowserUsername,
			ResourceUrl:         sharePost.RemotePath,
			DownloadPage:        sharePost.DownloadPage,
			DownloadUrl:         sharePost.DownloadUrl,
		}
	} else {
		wd_log.Debugf("=> share pwd: %s", sharePost.DownloadPasswd)
		wd_log.Info("=> share with password")
		shareFileBrowserUpload = wd_share_file_browser_upload.WdShareFileBrowserUpload{
			IsSendSuccess:       true,
			HostUrl:             p.Settings.FileBrowserBaseConfig.FileBrowserHost,
			FileBrowserUserName: p.Settings.FileBrowserBaseConfig.FileBrowserUsername,
			ResourceUrl:         sharePost.RemotePath,
			DownloadUrl:         sharePost.DownloadUrl,
			DownloadPage:        sharePost.DownloadPage,
			DownloadPasswd:      sharePost.DownloadPasswd,
		}
	}
	wd_log.Infof("=> share user name: %s", p.Settings.FileBrowserBaseConfig.FileBrowserUsername)
	wd_log.Infof("=> share remote path: %s", sharePost.RemotePath)

	p.shareFileBrowserUpload = &shareFileBrowserUpload
	wd_log.DebugJsonf(p.shareFileBrowserUpload, "shareFileBrowserUpload changes by send seccess\n")
	return nil
}

func genPwd(cnt uint) string {
	if cnt == 0 {
		return ""
	}

	return randomStrBySed(cnt, randPasswordSeed)
}

// randomStr
// new random string by cnt
//
//nolint:golint,unused
func randomStrBySed(cnt uint, sed string) string {
	var letters = []byte(sed)
	result := make([]byte, cnt)
	keyL := len(letters)
	rs := rand.New(rand.NewSource(time.Now().Unix()))
	for i := range result {
		result[i] = letters[rs.Intn(keyL)]
	}
	return string(result)
}

func (p *FileBrowserPlugin) saveStepsTransfer() error {
	if p.Settings.StepsOutDisable {
		wd_log.Debugf("steps out disable by flag [ %v ], skip save steps transfer", p.Settings.StepsOutDisable)
		return nil
	}

	if p.shareFileBrowserUpload != nil {
		_, errStepsTransfer := wd_steps_transfer.Out(p.Settings.RootPath, p.Settings.StepsTransferPath, *p.woodpeckerInfo,
			wd_share_file_browser_upload.WdShareKeyFileBrowserUpload, *p.shareFileBrowserUpload)
		if errStepsTransfer != nil {
			return errStepsTransfer
		}
	}

	return nil
}
