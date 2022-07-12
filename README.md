# Docker Compose v2 +

A docker compose enhanced tool. 

Additional support: 

- **HOOKs**

  executing shell, command, [golang script](https://github.com/goplus/gop
  )(via [interpreter](https://github.com/goplus/igop))

- Copy file/folder from the image to the local filesystem.


> Base on [docker/compose v2.6.1](https://github.com/docker/compose), Follow official updates unscheduled.

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

Create and start containers with HOOKs, be used in place of `docker compose up`.

- `[OPTIONS]`: the options of `docker compose --help`
- `[SERVICE...]`: the list of services that you want to `up`
- `[OPTIONS_OF_UP]`: the options of `docker compose up --help`
- `--pull` (default: false): pull the image before `up`
- `--hook` (default: false): executing commands before/after `up`

> You can specify any custom arguments(see [cli example](#CLI)), they can be read in the shell/golang scripts

- **pre-deploy**: Array of command, executing before `up`
- **post-deploy**: Array of command, executing after `up`

#### Execution sequence

1. `docker compose pull` if `--pull` be specified
2. **pre-deploy** of hooks
   1. command 1
   2. command 2
   3. ...
3. `docker compose up`
4. **post-deploy** of hooks
    1. command 1
    2. command 2
    3. ...

#### Examples

Files were in `/this/project/examples/`, copy to `/a/b/` where you want to put

```
- /a/b/
    - docker-compose.yaml
    - conf/
      - app.json
    - scripts/
      - main.go
      - def.go
      - abc.sh
```

##### docker-compose.yaml

```
x-hooks:
  pre-deploy:
    - []
    - []
  post-deploy:
    - []
    - []

services:
  ...
```

##### CLI

like `docker compose up`

```
docker compose -f '/a/b/docker-compose.yaml' deploy service-1 service-2 -d --hook --other-arg1 --other-arg2
```
or
```
cd /a/b/
docker compose deploy service-1 service-2 -d --hook --other-arg1 --other-arg2
```

#### Command specs:

- **command**: any command like `["echo", "\"hello\""]`
- **shell-key**: executing an inline shell of key starts with "x-"
```
["shell-key", "x-a-b-shell"]
```
- **igo-key**: executing an inline [gop script](https://goplus.org/) of key starts with "x-"
```
["igo-key", "x-b-c-igo"]
```
- **igo-path**: a path of golang file included `package main` & `func main()`
```
["igo-path", "scripts/main.go"]
```

#### Relative path/working directory

1. All path in the `pre-deploy/post-deploy` are relative to the `docker-compose.yaml` if you set a relative path, eg: `scripts/main.go` is `/a/b/scripts/main.go`

2. Working directory is the directory of `docker-compose.yaml`, eg: `/a/b/`

#### Execution arguments

- **command**: nothing will change

```
$ cd /a/b
$ echo "hello"
```

- **shell-key**: includes all arguments

```
$ cd /a/b
$ /usr/bin/sh /a/b/x-a-b-shell.sh -f '/a/b/docker-compose.yaml' deploy service-1 service-2 -d --hook --other-arg1 --other-arg2
```

- **igo-key**: includes all arguments
 
```
$ cd /a/b 
$ /a/b/x-b-c-igo.gop -f '/a/b/docker-compose.yaml' deploy service-1 service-2 -d --hook --other-arg1 --other-arg2
```

- **igo-path**: includes all arguments

```
$ cd /a/b 
$ /a/b/scripts/main.go -f '/a/b/docker-compose.yaml' deploy service-1 service-2 -d --hook --other-arg1 --other-arg2
```

### Golang script

ToDo
