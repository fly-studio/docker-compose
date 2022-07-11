# Docker Compose with hooks v2 

A docker compose tool with hook, supported executing shell, command, [gop](https://github.com/goplus/gop
)(golang script, [interpreter](https://github.com/goplus/igop))

## Features

### docker compose deploy --pull --hook

`deploy` is as same as `docker compose up`, but added the following parameters

- `--pull` (default: false): pull the image before `up`
- `--hook` (default: false): executing some command before/after `up`

> Warning: 
> You must specify the parameters explicitly, or they will not be executed

### Execution sequence

- **pull**, as same as `docker compose pull`
- **pre-deploy** of x-hooks
    - command 1
    - command 2
    - ...
- **up**, as same as `docker compose up`
- **post-deploy** of x-hooks
    - command 1
    - command 2
    - ...

## Hooks

Executing some command/shell/go scripts before/after `up`

- pre-deploy/post-deploy: Array
### Supported:
- **command**: like `["echo", "\"hello\""]`
- **shell-key**: executing an inline shell with "x-" name, look "x-a-b-shell"
- **igo-key**: executing an inline [gop script](https://goplus.org/) with "x-" name, look "x-b-c-igo"
- **igo-path**: a path of golang file included `package main` & `func main()`

### Examples

Files was in `/this/project/examples/`, copy to `/a/b/` where you want to put it

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

CLI

```
dock compose -f '/a/b/docker-compose.yaml' deploy service-1 service-2 -d --hook --other-arg1 --other-arg2
```

or 
```
cd /a/b/
dock compose deploy service-1 service-2 -d --hook --other-arg1 --other-arg2
```

### Relative path/working directory

1. All path in the `pre-deploy/post-deploy` are relative to the `docker-compose.yaml` if you set a relative path, eg: `/a/b/` 
2. Working directory of command is the path of `docker-compose.yaml`, eg: `/a/b/`

### Execution arguments

- **command**: Nothing will change

- **shell-key**: Append all arguments

```
cd /a/b && /usr/bin/sh x-a-b-shell.sh -f '/a/b/docker-compose.yaml' deploy service-1 service-2 -d --hook --other-arg1 --other-arg2
```

- **igo-key**: Append all arguments
 
```
cd /a/b && ./x-b-c-igo.gop -f '/a/b/docker-compose.yaml' deploy service-1 service-2 -d --hook --other-arg1 --other-arg2
```

- **igo-path**: Append all arguments

```
cd /a/b && ./scripts/main.go -f '/a/b/docker-compose.yaml' deploy service-1 service-2 -d --hook --other-arg1 --other-arg2
```

## Golang script

ToDo
