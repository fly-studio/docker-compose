# Docker Compose with hooks v2 

A docker compose tool with hooks, supported executing shell, command, [gop](https://github.com/goplus/gop
)(golang script, [interpreter](https://github.com/goplus/igop))

Base on [docker/compose v2.6.1](https://github.com/docker/compose), Follow official updates unscheduled.

## Install

Copy the [release](https://github.com/fly-studio/docker-compose/releases) to 
```
/usr/libexec/docker/cli-plugins/docker-compose 
```

or (ln -s from above is recommended)

```
/usr/bin/docker-compose
```

## Features

### docker compose deploy --pull --hook

`deploy` is as same as `docker compose up`, but added the following arguments

- `--pull` (default: false): pull the image before `up`
- `--hook` (default: false): executing some command before/after `up`

### Execution sequence

- `docker compose pull`
- **pre-deploy** of x-hooks
    - command 1
    - command 2
    - ...
- `docker compose up`
- **post-deploy** of x-hooks
    - command 1
    - command 2
    - ...

## Examples

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

docker-compose.yaml:
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

## Usage

like `docker compose up`

```
docker compose -f '/a/b/docker-compose.yaml' deploy service-1 service-2 -d --hook --other-arg1 --other-arg2
```

or
```
cd /a/b/
docker compose deploy service-1 service-2 -d --hook --other-arg1 --other-arg2
```

You can set custom arguments, they can be read in the shell/golang scripts

## Hooks

Executing some command/shell/go scripts before/after `up`


- **pre-deploy**: Array of command
- **post-deploy**: Array of command

### Command specs:

- **command**: like `["echo", "\"hello\""]`
- **shell-key**: executing an inline shell of name starts with "x-"
```
- ["shell-key", "x-a-b-shell"]
```
- **igo-key**: executing an inline [gop script](https://goplus.org/) of name starts with "x-" name
```
- ["igo-key", "x-b-c-igo"]
```
- **igo-path**: a path of golang file included `package main` & `func main()`
```
- ["igo-path", "scripts/main.go"]
```

### Relative path/working directory

1. All path in the `pre-deploy/post-deploy` are relative to the `docker-compose.yaml` if you set a relative path, eg: `/a/b/` 
2. Working directory of command is the path of `docker-compose.yaml`, eg: `/a/b/`

### Execution arguments

- **command**: nothing will change

```
cd /a/b
echo "hello"
```

- **shell-key**: all arguments

```
cd /a/b
/usr/bin/sh x-a-b-shell.sh -f '/a/b/docker-compose.yaml' deploy service-1 service-2 -d --hook --other-arg1 --other-arg2
```

- **igo-key**: all arguments
 
```
cd /a/b 
./x-b-c-igo.gop -f '/a/b/docker-compose.yaml' deploy service-1 service-2 -d --hook --other-arg1 --other-arg2
```

- **igo-path**: all arguments

```
cd /a/b 
/scripts/main.go -f '/a/b/docker-compose.yaml' deploy service-1 service-2 -d --hook --other-arg1 --other-arg2
```

## Golang script

ToDo
