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
- [Golang script](#Golang script)

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

> You can specify any custom arguments(see [cli example](#Quick start)), they can be read in the shell/golang scripts

- **pre-deploy**: Array of command, executing before `up`
- **post-deploy**: Array of command, executing after `up`

#### Execution sequence

1. `docker compose pull [SERVICE...]` if `--pull` be specified
2. pre-deploy of global
3. pre-deploy of each service of `[SERVICE...]` 
4. `docker compose up [SERVICE...]`
5. post-deploy of each service of `[SERVICE...]`
6. post-deploy of global

#### Relative path/working directory

1. All path in the `pre-deploy/post-deploy` are relative to the `docker-compose.yaml` if you set a relative path, eg: `scripts/main.go` is `/a/b/scripts/main.go`
2. Working directory is the directory of `docker-compose.yaml`, eg: `/a/b/`

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
    - ["shell-key", "x-a-b-shell"]
  post-deploy:
    - []
x-a-b-shell: |
  echo "hello"
      
services:
  nginx:
    container_name: nginx
    image: nginx:latest
    x-hooks:
      pre-deploy:
        - ["shell-key", "x-a-b-shell", "3"]
      post-deploy:
        - []
    x-a-b-shell: |
      ping 8.8.8.8 -c $2
```

##### Quick start

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

##### · command

Any command 

```
- ["echo", "\"hello\""]`
```

##### · shell-key

executing an inline shell of key starts with "x-".

- global pre-deploy/post-deploy will find the root key of "x-"
- service pre-deploy/post-deploy, the key in service is preferred, then the root

```
- ["shell-key", "x-a-b-shell"]
```
Go to
```
$ cd /a/b
$ /usr/bin/sh /a/b/x-a-b-shell.sh
```

##### · igo-key

executing an inline [gop script](https://goplus.org/) of key starts with "x-"

- global pre-deploy/post-deploy will find the root key of "x-"
- service pre-deploy/post-deploy, the key in service is preferred, then the root


```
- ["igo-key", "x-b-c-igo"]
```
Go to
```
$ cd /a/b
$ igop /a/b/x-b-c-igo.gop
```

##### · igo-path

a path of golang file included `package main` & `func main()`

```
["igo-path", "scripts/main.go"]
```
Go to
```
$ cd /a/b 
$ igop /a/b/scripts/main.go 
```

#### Execution arguments

##### · No change
```
- ["echo", "'hello' > 1.txt"]
- ["igo-key", "x-b-c-igo", "--other"]

$ echo 'hello' > 1.txt
$ igop /a/b/x-b-c-igo.gop --other
```

##### · Environment

You can use the environment in `.env` or exported
```
$ export SERVER_ID = 2
- ["shell-key", "x-a-b-shell", "--server-id", "${SERVER_ID}"]

$ sh /a/b/x-a-b-shell.sh --server-id 2
```

##### · The arguments of `docker compose deploy ....`

append all arguments

```
- ["igo-path", "scripts/main.go", "--custom", "1", "{ARGS}"]

$ igop /a/b/scripts/main.go --custom 1  -f '/a/b/docker-compose.yaml' deploy service-1 service-2 -d --hook --other-arg1 --other-arg2
```


### Golang script

ToDo
