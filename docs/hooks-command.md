# Command of Hooks 

## Examples

Files were in `/this/project/examples/`, copy to `/a/b/` where you want to put

```
- /a/b/
    - docker-compose.yaml
    - scripts/
      - main.go
      - def.go
      - abc.sh
```

### docker-compose.yaml

```
x-hooks:
  pre-deploy:
    - ["shell-key", "x-a-b-shell"]
x-a-b-shell: |
  echo "hello"
      
services:
  nginx:
    container_name: nginx
    image: nginx:latest
    volumns:
        - /local/nginx/conf:/etc/nginx/conf
    x-hooks:
      pre-deploy:.
        - ["docker", "compose", "cpi", "nginx", "/etc/nginx/conf", "/local/nginx/"]
        - ["shell-key", "x-a-b-shell", "3"]
    x-a-b-shell: |
      ping 8.8.8.8 -c $2
  redis:
    container_name: redis
    image: redis:latest
    
```

### Quick start

like `docker compose up`

```
docker compose -f '/a/b/docker-compose.yaml' deploy nginx redis -d --pull --hook --other-arg1 --other-arg2
```

or

```
cd /a/b/
docker compose deploy nginx redis -d --hook --other-arg1 --other-arg2
```

## Command specs:

### · command

Any command

```
- ["echo", "\"hello\""]
```

### · shell-key

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

### · igo-key

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

### · igo-path

a path of golang file included `package main` & `func main()`

```
["igo-path", "scripts/main.go"]
```
Go to
```
$ cd /a/b 
$ igop /a/b/scripts/main.go 
```

## Command arguments

### · No change
```
- ["echo", "'hello' > 1.txt"]
- ["igo-key", "x-b-c-igo", "--other"]

$ echo 'hello' > 1.txt
$ igop /a/b/x-b-c-igo.gop --other
```

### · Environment

You can use the environment from `.env` or exported
```
$ export SERVER_ID = 2
- ["shell-key", "x-a-b-shell", "--server-id", "${SERVER_ID}"]

$ sh /a/b/x-a-b-shell.sh --server-id 2
```

### · The arguments of `docker compose deploy ....`

append all arguments

```
- ["igo-path", "scripts/main.go", "--custom", "1", "{ARGS}"]

$ igop /a/b/scripts/main.go --custom 1  -f '/a/b/docker-compose.yaml' deploy service-1 service-2 -d --hook --other-arg1 --other-arg2
```

