package file_browser_upload

import (
	"fmt"
	"github.com/sinlov-go/go-common-lib/pkg/folder"
	"github.com/sinlov-go/go-common-lib/pkg/string_tools"
	"github.com/sinlov/filebrowser-client/file_browser_client"
	tools "github.com/sinlov/filebrowser-client/tools/str_tools"
	"github.com/sinlov/filebrowser-client/web_api"
	"github.com/woodpecker-kit/woodpecker-tools/wd_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_info_shot"
	"github.com/woodpecker-kit/woodpecker-tools/wd_log"
	"github.com/woodpecker-kit/woodpecker-tools/wd_template"
	"log"
	"math/rand"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type (
	// FileBrowserPlugin file_browser_upload all config
	FileBrowserPlugin struct {
		Name           string
		Version        string
		WoodpeckerInfo *wd_info.WoodpeckerInfo
		Config         Config

		FuncPlugin FuncPlugin `json:"-"`
	}
)

type FuncPlugin interface {
	Exec() error

	loadStepsTransfer() error
	checkArgs() error
	saveStepsTransfer() error
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
	errCheck := argCheckInArr("args file_browser_dist_type", p.Config.FileBrowserSendConfig.FileBrowserDistType, pluginDistTypeSupport)
	if errCheck != nil {
		return errCheck
	}

	if len(p.Config.FileBrowserSendConfig.FileBrowserTargetFileGlob) == 0 && p.Config.FileBrowserSendConfig.FileBrowserTargetFileRegular == "" {
		return fmt.Errorf("args file_browser_target_file_glob and file_browser_target_file_regular not be empty")
	}

	if p.Config.FileBrowserBaseConfig.FileBrowserUsername == "" {
		return fmt.Errorf("args file_browser_username not be empty")
	}

	if p.Config.FileBrowserSendConfig.FileBrowserRemoteRootPath == "" {
		return fmt.Errorf("args file_browser_remote_root_path not be empty")
	}

	// check default FileBrowserTimeoutPushSecond
	if p.Config.FileBrowserBaseConfig.FileBrowserTimeoutPushSecond < 60 {
		p.Config.FileBrowserBaseConfig.FileBrowserTimeoutPushSecond = 60
	}

	// check default p.Config.FileBrowserBaseConfig.FileBrowserWorkSpace
	if p.Config.FileBrowserBaseConfig.FileBrowserWorkSpace == "" {
		p.Config.FileBrowserBaseConfig.FileBrowserWorkSpace = p.WoodpeckerInfo.BasicInfo.CIWorkspace
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

	wdInfoShort := wd_info_shot.ParseWoodpeckerInfo2Shot(*p.WoodpeckerInfo)

	fileBrowserClient, errNew := file_browser_client.NewClient(
		p.Config.FileBrowserBaseConfig.FileBrowserUsername,
		p.Config.FileBrowserBaseConfig.FileBrowserUserPassword,
		p.Config.FileBrowserBaseConfig.FileBrowserHost,
		p.Config.TimeoutSecond,
		p.Config.FileBrowserBaseConfig.FileBrowserTimeoutPushSecond,
	)

	if errNew != nil {
		return fmt.Errorf("new fileBrowser client err: %v", errNew)
	}
	fileBrowserClient.Debug(p.Config.Debug)

	var remoteRealRootPath = strings.TrimRight(p.Config.FileBrowserSendConfig.FileBrowserRemoteRootPath, "/")

	switch p.Config.FileBrowserSendConfig.FileBrowserDistType {
	default:
		return fmt.Errorf("send dist type not support %s", p.Config.FileBrowserSendConfig.FileBrowserDistType)
	case DistTypeGit:
		commitShortSha := string([]rune(p.WoodpeckerInfo.CurrentInfo.CurrentCommitInfo.CiCommitSha))[:8]
		if p.WoodpeckerInfo.CurrentInfo.CurrentCommitInfo.CiCommitTag != "" {

			tagPath, errPathTag := wd_template.RenderTrim(distGitGraphTag, wdInfoShort)
			if errPathTag != nil {
				return fmt.Errorf("render as %s \nerr: %v", distGitGraphTag, errPathTag)
			}
			remoteRealRootPath = path.Join(remoteRealRootPath, tagPath, commitShortSha)
		} else if p.WoodpeckerInfo.CurrentInfo.CurrentCommitInfo.CiCommitPullRequest != "" {
			prPath, errPathPr := wd_template.RenderTrim(distGitGraphPullRequest, wdInfoShort)
			if errPathPr != nil {
				return fmt.Errorf("render as %s \nerr: %v", distGitGraphPullRequest, errPathPr)
			}
			remoteRealRootPath = path.Join(remoteRealRootPath, prPath, commitShortSha)
		} else {
			defaultPath, errPathDefault := wd_template.RenderTrim(distGitGraphDefault, wdInfoShort)
			if errPathDefault != nil {
				return fmt.Errorf("render as %s \nerr: %v", distGitGraphDefault, errPathDefault)
			}

			remoteRealRootPath = path.Join(remoteRealRootPath, defaultPath, commitShortSha)
		}
	case DistTypeCustom:
		renderPath, err := wd_template.RenderTrim(p.Config.FileBrowserSendConfig.FileBrowserDistGraph, wdInfoShort)
		if err != nil {
			return fmt.Errorf("setting file_browser_dist_graph as %s \nerr: %v", p.Config.FileBrowserSendConfig.FileBrowserDistGraph, err)
		}
		remoteRealRootPath = path.Join(remoteRealRootPath, renderPath)
	}

	wd_log.Debugf("final remoteRealRootPath: %s", remoteRealRootPath)

	targetRootPath := filepath.Join(p.Config.FileBrowserBaseConfig.FileBrowserWorkSpace, p.Config.FileBrowserSendConfig.FileBrowserTargetDistRootPath)
	wd_log.Debugf("workOnSend fileBrowserDistType: %s", p.Config.FileBrowserSendConfig.FileBrowserDistType)
	if p.Config.FileBrowserSendConfig.FileBrowserDistType == DistTypeCustom {
		wd_log.Verbosef("workOnSend fileBrowserDistGraph: %s", p.Config.FileBrowserSendConfig.FileBrowserDistGraph)
	}
	wd_log.Debugf("workOnSend remoteRealRootPath: %s", remoteRealRootPath)
	wd_log.Debugf("workOnSend targetRootPath: %s", targetRootPath)
	wd_log.Debugf("workOnSend fileBrowserWorkSpace: %s", p.Config.FileBrowserBaseConfig.FileBrowserWorkSpace)
	wd_log.Debugf("workOnSend targetDistRootPath: %s", p.Config.FileBrowserSendConfig.FileBrowserTargetDistRootPath)
	wd_log.Debugf("workOnSend target File Regular: %v", p.Config.FileBrowserSendConfig.FileBrowserTargetFileRegular)
	wd_log.Debugf("workOnSend target File Glob: %v", p.Config.FileBrowserSendConfig.FileBrowserTargetFileGlob)

	if !(folder.PathExistsFast(targetRootPath)) {
		return fmt.Errorf("file browser want send file local path not exists at: %s", targetRootPath)
	}

	var fileSendPathList []string
	if folder.PathIsFile(targetRootPath) {
		wd_log.Debugf("target path is file just send one local file: %s", targetRootPath)
		fileSendPathList = append(fileSendPathList, targetRootPath)
	} else {
		if p.Config.FileBrowserSendConfig.FileBrowserTargetFileGlob != nil && len(p.Config.FileBrowserSendConfig.FileBrowserTargetFileGlob) > 0 {
			wd_log.Debugf("target file want find by File Glob: %v", p.Config.FileBrowserSendConfig.FileBrowserTargetFileGlob)
			for _, glob := range p.Config.FileBrowserSendConfig.FileBrowserTargetFileGlob {
				walkByGlob, errWalkAllByGlob := folder.WalkAllByGlob(targetRootPath, glob, true)
				if errWalkAllByGlob != nil {
					return fmt.Errorf("file browser want send file local path with glob %s be err: %v", targetRootPath, errWalkAllByGlob)
				}
				wd_log.Debugf("target path find by File Glob [ %s ] files:\n%s", glob, strings.Join(walkByGlob, "\n"))
				fileSendPathList = append(fileSendPathList, walkByGlob...)
			}
		}
		if p.Config.FileBrowserSendConfig.FileBrowserTargetFileRegular != "" {
			wd_log.Debugf("target file want find by File Regular: %v", p.Config.FileBrowserSendConfig.FileBrowserTargetFileRegular)

			matchPath, err := folder.WalkAllByMatchPath(targetRootPath, p.Config.FileBrowserSendConfig.FileBrowserTargetFileRegular, true)
			if err != nil {
				return fmt.Errorf("file browser want send file local path with file regular %s be err: %v", targetRootPath, err)
			}
			wd_log.Debugf("target path find by File Regular [ %s ] files:\n%s", p.Config.FileBrowserSendConfig.FileBrowserTargetFileRegular, strings.Join(matchPath, "\n"))
			fileSendPathList = append(fileSendPathList, matchPath...)
		}
	}

	if len(fileSendPathList) == 0 {
		return fmt.Errorf("file browser want send file local path not find any file at path: %s\nglob: %s\nregular: %s",
			targetRootPath,
			p.Config.FileBrowserSendConfig.FileBrowserTargetFileGlob,
			p.Config.FileBrowserSendConfig.FileBrowserTargetFileRegular,
		)
	}

	wd_log.Debugf("now send path len %d", len(fileSendPathList))
	fileSendPathList = tools.StrArrRemoveDuplicates(fileSendPathList)
	if p.Config.Debug {
		wd_log.Debugf("debug: send path remove duplicates len %d", len(fileSendPathList))
	}

	if p.Config.DryRun {
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
		errSendOneFile := fileBrowserClient.ResourcesPostFile(resourcePostOne, p.Config.Debug)
		if err != nil {
			return errSendOneFile
		}
		if p.Config.FileBrowserSendConfig.FileBrowserShareLinkEnable {
			errSendFileShare := shareBySendConfig(fileBrowserClient, *p, remotePath, false)
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
			errSendOneFile := fileBrowserClient.ResourcesPostFile(resourcePost, p.Config.Debug)
			if err != nil {
				return errSendOneFile
			}
		}
		if p.Config.FileBrowserSendConfig.FileBrowserShareLinkEnable {
			errSendFileShare := shareBySendConfig(fileBrowserClient, *p, remoteRealRootPath, true)
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

func shareBySendConfig(client file_browser_client.FileBrowserClient, p FileBrowserPlugin, remotePath string, isDir bool) error {
	expires := strconv.Itoa(int(p.Config.FileBrowserSendConfig.FileBrowserShareLinkExpires))
	passWord := p.Config.FileBrowserSendConfig.FileBrowserShareLinkPassword
	if p.Config.FileBrowserSendConfig.FileBrowserShareLinkAutoPasswordEnable {
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
			Unit:     p.Config.FileBrowserSendConfig.FileBrowserShareLinkUnit,
		},
	}
	sharePost, errSendShareFile := client.SharePost(shareResource)
	if errSendShareFile != nil {
		return errSendShareFile
	}
	log.Printf("=> share page: %s", sharePost.DownloadPage)
	if passWord != "" {
		log.Printf("=> share pwd: %s", sharePost.DownloadPasswd)
	}
	log.Printf("=> share user name: %s", p.Config.FileBrowserBaseConfig.FileBrowserUsername)
	log.Printf("=> share remote path: %s", sharePost.RemotePath)
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

	return nil
}
