# Docker Compose v2 +

A docker compose enhanced tool. 

> Base on [docker/compose v2.6.1](https://github.com/docker/compose), Follow official updates unscheduled.

---

Additional support: 

- HOOKs

  executing shell, command, [golang script](https://github.com/goplus/gop
  )(via [interpreter](https://github.com/goplus/igop))

- Copy file/folder from the image to the local filesystem.


## TOC

- [Install](#Install)
- [Usage](#Usage)
  - [Copy from image](#Copy-from-image)
  - [Hooks](#Hooks)
- [Golang script](#Golang-script)

## Install

Copy the [release](https://github.com/fly-studio/docker-compose/releases) to 
```
/usr/libexec/docker/cli-plugins/docker-compose 
```

or (`ln -s` is recommended)

```
/usr/bin/docker-compose
```

## Usage

### Copy from image

```
docker compose [OPTIONS] cpi [SERVICE] [PATH_IN_IMAGE] [LOCAL_PATH] --follow-link
```

Copy a file/folder from the image of the [SERVICE] to the local filesystem

- `[OPTIONS]`: the options of `docker compose --help`
- `[SERVICE]`: the service name that you want to copy from
- `[PATH_IN_IMAGE]`: the path in the image of the `[SERVICE]`, source path
- `[LOCAL_PATH]`: the path of local filesystem, destination path
- `--follow-link | -L`: always follow symbol link in `[PATH_IN_IMAGE]`

#### Examples

```
docker compose -f "/a/b/docker-compose.yaml" cpi nginx /etc/nginx/conf /local/nginx-conf/
```

### Hooks

```
docker compose [OPTIONS] deploy [SERVICE...] [OPTIONS_OF_UP] --pull --hook
```

Create and start containers with HOOKs, the usage is similar to `docker compose up`.

- `[OPTIONS]`: the options of `docker compose --help`
- `[SERVICE...]`: the list of services that you want to `up`
- `[OPTIONS_OF_UP]`: the options of `docker compose up --help`
- `--pull` (default: false): pull the image before `up`
- `--hook` (default: false): executing commands before/after `up`

> You can specify any custom arguments(see [cli example](docs/hooks-command.md#Quick-start)), they can be read in the shell/golang scripts

#### Example
```
x-hooks:
  pre-deploy:
    - ["command", "args"]  
  post-deploy:
    - ["command", "args"]
    
services:
  nginx:
    image: nginx:latest
    x-hooks:
      pre-deploy:
        - ["command", "args"]  
      post-deploy:
        - ["command", "args"]
```

- **pre-deploy**: Array of command

- **post-deploy**: Array of command

#### Execution sequence

1. `docker compose pull [SERVICE...]` if `--pull` be specified
2. Global pre-deploy 
3. pre-deploy of each service of `[SERVICE...]`
4. `docker compose up [SERVICE...]`
5. post-deploy of each service of `[SERVICE...]`
6. Global post-deploy

#### Relative path/working directory

1. Path in the command is relative to the directory of `docker-compose.yaml`,  eg: `scripts/main.go` is `/a/b/scripts/main.go`

2. Working directory is the directory of `docker-compose.yaml`, eg: `/a/b/`

#### Command usage

  see [Command of Hooks](docs/hooks-command.md)