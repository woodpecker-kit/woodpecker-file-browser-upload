---
name: woodpecker-file-browser-upload
description: woodpecker file browser upload
author: woodpecker-kit
tags: [ file-browser ]
containerImage: sinlov/woodpecker-plugin-env
containerImageUrl: https://hub.docker.com/r/sinlov/woodpecker-plugin-env
url: https://github.com/woodpecker-kit/woodpecker-file-browser-upload
icon: https://woodpecker-ci.org/img/logo.svg
---

woodpecker-file-browser-upload

## Settings

| Name                                           | Required | Default value | Description                                                                                                     |
|------------------------------------------------|----------|---------------|-----------------------------------------------------------------------------------------------------------------|
| `debug`                                        | **no**   | *false*       | open debug log or open by env `PLUGIN_DEBUG`                                                                    |
| `timeout_second`                               | **no**   | *10*          | api timeout default: 10                                                                                         |
| `file_browser_timeout_push_second`             | **no**   | *60*          | push each file timeout push second, must gather than 60.default: 60                                             |
| `file_browser_host`                            | **yes**  | *none*        | file_browser host like http://127.0.0.1:80/                                                                     |
| `file_browser_username`                        | **yes**  | *none*        | file_browser username                                                                                           |
| `file_browser_user_password`                   | **yes**  | *none*        | file_browser user password                                                                                      |
| `file_browser_work_space`                      | **no**   | *none*        | file_browser work space. default "" will use env:CI_WORKSPACE                                                   |
| `file_browser_remote_root_path`                | **yes**  | *none*        | send to file_browser base path                                                                                  |
| `file_browser_dist_type`                       | **yes**  | *none*        | type of dist file graph only can use: git, custom                                                               |
| `file_browser_dist_graph`                      | **no**   | *none*        | type of dist custom                                                                                             |
| `file_browser_target_dist_root_path`           | **no**   | *""*          | path of file_browser work on root, can set "". default: ""                                                      |
| `file_browser_target_file_globs`               | **yes**  | *none*        | globs list of send to file_browser under file_browser_target_dist_root_path                                     |
| `file_browser_target_file_regular`             | **no**   | *none*        | regular of send to file_browser under file_browser_target_dist_root_path                                        |
| `file_browser_share_link_enable`               | **no**   | *true*        | share dist dir as link, default: true                                                                           |
| `file_browser_share_link_expires`              | **no**   | *0*           | if set 0, will allow share_link exist forever，default: 0                                                        |
| `file_browser_share_link_unit`                 | **no**   | *days*        | take effect by open share_link, only can use as `[ days hours minutes seconds ]`                                |
| `file_browser_share_link_password`             | **no**   | *""*          | password of share_link, if not set will not use password, default: ""                                           |
| `file_browser_share_link_auto_password_enable` | **no**   | *false*       | password of share_link auto , if open this will cover settings.file_browser_share_link_password. default: false |

**Hide Settings:**

| Name                                        | Required | Default value                     | Description                                                                      |
|---------------------------------------------|----------|-----------------------------------|----------------------------------------------------------------------------------|
| `timeout_second`                            | **no**   | *10*                              | "command timeout setting by second                                               |
| `woodpecker-kit-steps-transfer-file-path`   | **no**   | `*.woodpecker_kit.steps.transfer` | Steps transfer file path, default by `wd_steps_transfer.DefaultKitStepsFileName` |
| `woodpecker-kit-steps-transfer-disable-out` | **no**   | *false*                           | Steps transfer write disable out                                                 |

## Example

- workflow with backend `docker`

```yml
labels:
  backend: docker
steps:
  woodpecker-file-browser-upload:
    image: sinlov/woodpecker-file-browser-upload:latest
    pull: false
    settings:
      # debug: false # plugin debug switch
      file_browser_host: "http://127.0.0.1:80" # must set args, file_browser host like http://127.0.0.1:80
      file_browser_username: # must set args, file_browser username
        # https://woodpecker-ci.org/docs/usage/secrets
        from_secret: file_browser_user_name
      file_browser_user_password: # must set args, file_browser user password
        from_secret: file_browser_user_password
      file_browser_remote_root_path: dist/ # must set args, send to file_browser base path
      file_browser_target_file_globs: # must set args, globs list of send to file_browser under file_browser_target_dist_root_path
        - "**/*.tar.gz"
        - "**/*.sha256"
      file_browser_share_link_expires: 0 # if set 0, will allow share_link exist forever，default: 0
      file_browser_share_link_unit: days # take effect by open share_link, only can use as [ days hours minutes seconds ]
      file_browser_share_link_auto_password_enable: true # password of share_link auto , if open this will cover settings.file_browser_share_link_password. default: false
```

- workflow with backend `local`, must install at local and effective at evn `PATH`
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
      file_browser_host: "http://127.0.0.1:80" # must set args, file_browser host like http://127.0.0.1:80
      file_browser_username: # must set args, file_browser username
        # https://woodpecker-ci.org/docs/usage/secrets
        from_secret: file_browser_user_name
      file_browser_user_password: # must set args, file_browser user password
        from_secret: file_browser_user_password
      file_browser_remote_root_path: dist/ # must set args, send to file_browser base path
      file_browser_target_file_globs: # must set args, globs list of send to file_browser under file_browser_target_dist_root_path
        - "**/*.tar.gz"
        - "**/*.sha256"
      file_browser_share_link_expires: 0 # if set 0, will allow share_link exist forever，default: 0
      file_browser_share_link_unit: days # take effect by open share_link, only can use as [ days hours minutes seconds ]
      file_browser_share_link_auto_password_enable: true # password of share_link auto , if open this will cover settings.file_browser_share_link_password. default: false
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
      file_browser_timeout_push_second: 60 # push each file timeout push second, must gather than 60.default: 60
      file_browser_host: # must set args, file_browser host like http://127.0.0.1:80
        from_secret: file_browser_host
      file_browser_username: # must set args, file_browser username
        # https://woodpecker-ci.org/docs/usage/secrets
        from_secret: file_browser_user_name
      file_browser_user_password: # must set args, file_browser user password
        from_secret: file_browser_user_password
      file_browser_work_space: "" # file_browser work space. default "" will use env:CI_WORKSPACE
      file_browser_remote_root_path: dist/ # must set args, send to file_browser base path
      file_browser_dist_type: custom # must set args, type of dist file graph only can use: git, custom
      file_browser_dist_graph: "{{ Repo.HostName }}/{{ Repo.OwnerName }}/{{ Repo.ShortName }}/s/{{ Build.Number }}/{{ Build.Number }}-{{ Stage.Finished }}" # type of dist custom
      file_browser_target_dist_root_path: dist/ # path of file_browser work on root, can set "". default: ""
      file_browser_target_file_globs: # must set args, globs list of send to file_browser under file_browser_target_dist_root_path
        - "**/*.tar.gz"
        - "**/*.sha256"
      file_browser_target_file_regular: .*.tar.gz # must set args, regular of send to file_browser under file_browser_target_dist_root_path
      file_browser_share_link_enable: true # share dist dir as link, default: true
      file_browser_share_link_expires: 0 # if set 0, will allow share_link exist forever，default: 0
      file_browser_share_link_unit: days # take effect by open share_link, only can use as [ days hours minutes seconds ]
      file_browser_share_link_password: "" # password of share_link, if not set will not use password, default: ""
      file_browser_share_link_auto_password_enable: false # password of share_link auto , if open this will cover settings.file_browser_share_link_password. default: false
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
  "resource_url": "dist/woodpecker-kit/guidance-woodpecker-agent/tag/v1.0.0/10/9c764dd4/",
  "download_url": "http://192.168.50.199:59999/share/Q9yy3zSh",
  "download_passwd": "qsAbWK00"
}
```

### settings.debug

- if open `settings.debug` will try file browser use `override` for debug.
- if open `settings.woodpecker-kit-steps-transfer-disable-out` will disable out of `wd_steps_transfer`
- please close `settings.debug` in production models

### file_browser_dist_type

template use struct `wd_short_info.WoodpeckerInfoShort`

use file_browser_dist_type = `git`, send to filebrowser file tree like

```
# default
${file_browser_remote_root_path}/
	{{Repo.HostName}}/
		{{Repo.OwnerName}}/
			{{Repo.ShortName}}/
				b/
					{{Build.Number}}/
						{{Commit.Branch}}/
							{{Commit.Sha[0:8]}}

# if in pull request
${file_browser_remote_root_path}/
	{{Repo.HostName}}/
		{{Repo.OwnerName}}/
			{{Repo.ShortName}}/
				pr/
					{{Build.PR}}/
						{{Build.Number}}/
							{{Commit.Sha[0:8]}}

# if in tag
${file_browser_remote_root_path}/
	{{Repo.HostName}}/
		{{Repo.OwnerName}}/
			{{Repo.ShortName}}/
				tag/
					{{Build.Tag}}/
						{{Build.Number}}/
							{{Commit.Sha[0:8]}}
```

#### custom dist graph

- you can use file_browser_dist_type = `custom`, like

```
{{ Repo.HostName }}/{{ Repo.OwnerName }}/{{ Repo.ShortName }}/s/{{ Build.Number }}/{{ Build.Number }}-{{ Stage.Finished }}

// will out like this will append ${file_browser_remote_root_path}
dist/woodpecker-kit/guidance-woodpecker-agent/s/10/10-1705658166
```

