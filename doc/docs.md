---
name: woodpecker-file-browser-upload
description: woodpecker file browser upload
author: woodpecker-kit
tags: [ file-browser ]
containerImage: sinlov/woodpecker-plugin-env
containerImageUrl: https://hub.docker.com/r/sinlov/woodpecker-plugin-env
url: https://github.com/woodpecker-kit/woodpecker-file-browser-upload
icon: https://raw.githubusercontent.com/woodpecker-kit/woodpecker-file-browser-upload/main/doc/log_image.jpeg
---

woodpecker-file-browser-upload

## Settings

| Name                                      | Required | Default value | Description                                                                                                  |
|-------------------------------------------|----------|---------------|--------------------------------------------------------------------------------------------------------------|
| `debug`                                   | **no**   | *false*       | open debug log or open by env `PLUGIN_DEBUG`                                                                 |
| `file-browser-timeout-send-second`        | **no**   | *60*          | push each file timeout push second, must gather than 60.default: 60                                          |
| `file-browser-host`                       | **yes**  | *none*        | file_browser host like http://127.0.0.1:80                                                                   |
| `file-browser-username`                   | **yes**  | *none*        | file_browser username                                                                                        |
| `file-browser-user-password`              | **yes**  | *none*        | file_browser user password                                                                                   |
| `file-browser-work-space`                 | **no**   | *none*        | file_browser work space. default "" will use env:CI_WORKSPACE                                                |
| `file-browser-remote-root-path`           | **yes**  | *none*        | send to file_browser base path                                                                               |
| `file-browser-dist-type`                  | **yes**  | *none*        | type of dist file graph only can use: git, custom                                                            |
| `file-browser-dist-graph`                 | **no**   | *""*          | `file-browser-dist-type` setting `custom` dist graph define                                                  |
| `file-browser-target-dist-root-path`      | **no**   | *""*          | path of file_browser work on root, can set "". default: ""                                                   |
| `file-browser-file-glob`                  | **yes**  | *none*        | globs list of send to file_browser under file-browser-target-dist-root-path                                  |
| `file-browser-file-regular`               | **no**   | *none*        | regular of send to file_browser under file-browser-target-dist-root-path                                     |
| `file-browser-share-link-enable`          | **no**   | *false*       | share dist dir as link, default: false                                                                       |
| `file-browser-share-link-expire`          | **no**   | *0*           | if set 0, will allow share_link exist forever，default: 0                                                     |
| `file-browser-share-link-unit`            | **no**   | *days*        | take effect by open share_link, only can use as `[ days hours minutes seconds ]`                             |
| `file-browser-share-link-passwd`          | **no**   | *""*          | password of share_link, if not set will not use password, default: ""                                        |
| `file-browser-share-auto-password-enable` | **no**   | *false*       | password of share_link auto, if open this will cover settings.file-browser-share-link-passwd. default: false |

**Hide Settings:**

| Name                                        | Required | Default value                    | Description                                                                      |
|---------------------------------------------|----------|----------------------------------|----------------------------------------------------------------------------------|
| `timeout_second`                            | **no**   | *10*                             | command timeout setting by second                                                |
| `woodpecker-kit-steps-transfer-file-path`   | **no**   | `.woodpecker_kit.steps.transfer` | Steps transfer file path, default by `wd_steps_transfer.DefaultKitStepsFileName` |
| `woodpecker-kit-steps-transfer-disable-out` | **no**   | *false*                          | Steps transfer write disable out                                                 |

## Example

- workflow with backend `docker`

[![docker hub version semver](https://img.shields.io/docker/v/sinlov/woodpecker-file-browser-upload?sort=semver)](https://hub.docker.com/r/sinlov/woodpecker-file-browser-upload/tags?page=1&ordering=last_updated)
[![docker hub image size](https://img.shields.io/docker/image-size/sinlov/woodpecker-file-browser-upload)](https://hub.docker.com/r/sinlov/woodpecker-file-browser-upload)
[![docker hub image pulls](https://img.shields.io/docker/pulls/sinlov/woodpecker-file-browser-upload)](https://hub.docker.com/r/sinlov/woodpecker-file-browser-upload/tags?page=1&ordering=last_updated)

```yml
labels:
  backend: docker
steps:
  woodpecker-file-browser-upload:
    image: sinlov/woodpecker-file-browser-upload:latest
    pull: false
    settings:
      # debug: false # plugin debug switch
      file-browser-host: "http://127.0.0.1:80" # must set args, file_browser host like http://127.0.0.1:80
      file-browser-username: # must set args, file_browser username
        # https://woodpecker-ci.org/docs/usage/secrets
        from_secret: file_browser_user_name
      file-browser-user-password: # must set args, file_browser user password
        from_secret: file_browser_user_passwd
      file-browser-remote-root-path: dist/ # must set args, send to file_browser base path
      file-browser-dist-type: git # must set args, type of dist file graph only can use: git, custom
      file-browser-file-glob: # must set args, globs list of send to file_browser under file-browser-target-dist-root-path
        - "**/*.tar.gz"
        - "**/*.sha256"
      file-browser-share-link-enable: true # share dist dir as link, default: false
      file-browser-share-link-expire: 0 # if set 0, will allow share_link exist forever，default: 0
      file-browser-share-link-unit: days # take effect by open share_link, only can use as [ days hours minutes seconds ]
      file-browser-share-auto-password-enable: true # password of share_link auto, if open this will cover settings.file-browser-share-link-passwd. default: false
```

- workflow with backend `local`, must install at local and effective at evn `PATH`

[![GitHub license](https://img.shields.io/github/license/woodpecker-kit/woodpecker-file-browser-upload)](https://github.com/woodpecker-kit/woodpecker-file-browser-upload)
[![GitHub latest SemVer tag)](https://img.shields.io/github/v/tag/woodpecker-kit/woodpecker-file-browser-upload)](https://github.com/woodpecker-kit/woodpecker-file-browser-upload/tags)
[![GitHub release)](https://img.shields.io/github/v/release/woodpecker-kit/woodpecker-file-browser-upload)](https://github.com/woodpecker-kit/woodpecker-file-browser-upload/releases)

- install at ${GOPATH}/bin, latest

```bash
go install -a github.com/woodpecker-kit/woodpecker-file-browser-upload/cmd/woodpecker-file-browser-upload@latest
```

- install at ${GOPATH}/bin, v1.0.0

```bash
go install -v github.com/woodpecker-kit/woodpecker-file-browser-upload/cmd/woodpecker-file-browser-upload@v1.0.0
```

```yml
labels:
  backend: local
steps:
  woodpecker-file-browser-upload:
    image: woodpecker-file-browser-upload
    settings:
      # debug: false # plugin debug switch
      file-browser-host: "http://127.0.0.1:80" # must set args, file_browser host like http://127.0.0.1:80
      file-browser-username: # must set args, file_browser username
        # https://woodpecker-ci.org/docs/usage/secrets
        from_secret: file_browser_user_name
      file-browser-user-password: # must set args, file_browser user password
        from_secret: file_browser_user_passwd
      file-browser-remote-root-path: dist/ # must set args, send to file_browser base path
      file-browser-dist-type: git # must set args, type of dist file graph only can use: git, custom
      file-browser-file-glob: # must set args, globs list of send to file_browser under file-browser-target-dist-root-path
        - "**/*.tar.gz"
        - "**/*.sha256"
      file-browser-share-link-enable: true # share dist dir as link, default: false
      file-browser-share-link-expire: 0 # if set 0, will allow share_link exist forever，default: 0
      file-browser-share-link-unit: days # take effect by open share_link, only can use as [ days hours minutes seconds ]
      file-browser-share-auto-password-enable: true # password of share_link auto, if open this will cover settings.file-browser-share-link-passwd. default: false
```

- full config

```yaml
labels:
  backend: docker
steps:
  woodpecker-file-browser-upload:
    image: sinlov/woodpecker-file-browser-upload:latest
    pull: false
    settings:
      debug: false # plugin debug switch
      timeout_second: 10 # api timeout default: 10
      file-browser-timeout-send-second: 60 # push each file timeout push second, must gather than 60.default: 60
      file-browser-host: # must set args, file_browser base url, like http://127.0.0.1:80
        # https://woodpecker-ci.org/docs/usage/secrets
        from_secret: file_browser_host_url
      file-browser-username: # must set args, file_browser username
        from_secret: file_browser_user_name
      file-browser-user-password: # must set args, file_browser user password
        from_secret: file_browser_user_passwd
      file-browser-work-space: "" # file_browser work space. default "" will use env:CI_WORKSPACE
      file-browser-remote-root-path: dist/ # must set args, send to file_browser base path
      file-browser-dist-type: custom # must set args, type of dist file graph only can use: git, custom
      file-browser-dist-graph: "{{ Repo.Hostname }}/{{ Repo.OwnerName }}/{{ Repo.ShortName }}/s/{{ Build.Number }}/{{ Build.Number }}-{{ Stage.Finished }}"
      file-browser-target-dist-root-path: dist/ # path of file_browser work on root, can set "". default: ""
      file-browser-file-glob: # must set args, globs list of send to file_browser under file-browser-target-dist-root-path
        - "**/*.tar.gz"
        - "**/*.sha256"
      file-browser-file-regular: .*.tar.gz # must set args, regular of send to file_browser under file-browser-target-dist-root-path
      file-browser-share-link-enable: true # share dist dir as link, default: false
      file-browser-share-link-expire: 0 # if set 0, will allow share_link exist forever，default: 0
      file-browser-share-link-unit: days # take effect by open share_link, only can use as [ days hours minutes seconds ]
      file-browser-share-link-passwd: "" # password of share_link, if not set will not use password, default: ""
      file-browser-share-auto-password-enable: false # password of share_link auto , if open this will cover settings.file-browser-share-link-passwd. default: false
```

#### Out

- just add `.woodpecker_kit.steps.transfer` at git ignore
- will add out key `wd_share_file_browser_upload.WdShareKeyFileBrowserUpload` struct
  as `wd_share_file_browser_upload.WdShareFileBrowserUpload`

```json
{
  "is_send_success": true,
  "host_url": "http://192.168.50.199:59999/",
  "file_browser_user_name": "share",
  "download_url": "http://192.168.50.54:59999/api/public/dl/bIf-VmIz?token=HH8LYsnNi3LRpdpeFsYGp5HaUUjC8A4BHOOVsVNlmLq5VhpW7-lD1iTUKlcydbwb9zRp6GD--0eafpJMeI40i2MZX8-VFhJ0RJ_TFTZgHn2qLgfwFC9iwW2S1dlgA6c_",
  "download_page": "http://192.168.50.54:59999/share/bIf-VmIz",
  "download_passwd": "hfBhZjEC"
}
```

### settings.debug

- if open `settings.debug` will try file browser use `override` for debug.
- if open `settings.woodpecker-kit-steps-transfer-disable-out` will disable out of `wd_steps_transfer`
- please close `settings.debug` in production models

### file-browser-dist-type

template use struct `wd_short_info.WoodpeckerInfoShort`

use file-browser-dist-type = `git`, send to filebrowser file tree like

```
# default
${file-browser-remote-root-path}/
	{{Repo.Hostname}}/
		{{Repo.OwnerName}}/
			{{Repo.ShortName}}/
				b/
					{{Build.Number}}/
						{{Commit.Branch}}/
							{{Commit.Sha[0:8]}}

# if in pull request
${file-browser-remote-root-path}/
	{{Repo.Hostname}}/
		{{Repo.OwnerName}}/
			{{Repo.ShortName}}/
				pr/
					{{Build.PR}}/
						{{Build.Number}}/
							{{Commit.Sha[0:8]}}

# if in tag
${file-browser-remote-root-path}/
	{{Repo.Hostname}}/
		{{Repo.OwnerName}}/
			{{Repo.ShortName}}/
				tag/
					{{Build.Tag}}/
						{{Build.Number}}/
							{{Commit.Sha[0:8]}}
```

#### custom dist graph

- you can use file-browser-dist-type = `custom`, like

```
{{ Repo.Hostname }}/{{ Repo.OwnerName }}/{{ Repo.ShortName }}/s/{{ Build.Number }}/{{ Build.Number }}-{{ Stage.Finished }}

// will out like this will append ${file-browser-remote-root-path} as: dist/
dist/gitea.domain.com/woodpecker-kit/guidance-woodpecker-agent/s/10/10-1705658166
```

