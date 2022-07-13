package compose

import (
	"bytes"
	"context"
	"fmt"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/compose/v2/igo"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/sanathkr/go-yaml"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type xHooks struct {
	PreDeploy  []types.ShellCommand `mapstructure:"pre-deploy"`
	PostDeploy []types.ShellCommand `mapstructure:"post-deploy"`
}

type hook struct {
	ctx       context.Context
	cmd       *cobra.Command
	project   *types.Project
	backend   api.Service
	dockerCli command.Cli

	globalHooks   xHooks
	servicesHooks map[string]xHooks
}

func newHook(ctx context.Context, cmd *cobra.Command, dockerCli command.Cli, backend api.Service, project *types.Project) *hook {
	return &hook{
		ctx:           ctx,
		cmd:           cmd,
		project:       project,
		backend:       backend,
		dockerCli:     dockerCli,
		globalHooks:   xHooks{},
		servicesHooks: map[string]xHooks{},
	}
}

func (h *hook) parse() error {
	if err := loader.Transform(h.project.Extensions["x-hooks"], &h.globalHooks); err != nil {
		return err
	}

	for _, service := range h.project.Services {
		if _, ok := service.Extensions["x-hooks"]; ok {
			var xHooks xHooks
			if err := loader.Transform(service.Extensions["x-hooks"], &xHooks); err != nil {
				return err
			}
			h.servicesHooks[service.Name] = xHooks
		}
	}

	return nil
}

func (h *hook) PreDeploy(createOptions createOptions, upOptions upOptions, pullOptions pullOptions, service *types.ServiceConfig) error {
	xHook := h.globalHooks
	if service != nil {
		var ok bool
		if xHook, ok = h.servicesHooks[service.Name]; !ok { // service中不存在 x-hook则退出
			return nil
		}
		fmt.Printf("[Hook]pre-deploy of service: \"%s\"\n", service.Name)
	} else {
		fmt.Printf("[Hook]pre-deploy of global\n")
	}

	return h.handle(xHook.PreDeploy, createOptions, upOptions, pullOptions, service)
}

func (h *hook) PostDeploy(createOptions createOptions, upOptions upOptions, pullOptions pullOptions, service *types.ServiceConfig) error {
	xHook := h.globalHooks
	if service != nil {
		var ok bool
		if xHook, ok = h.servicesHooks[service.Name]; !ok { // service中不存在 x-hook则退出
			return nil
		}
		fmt.Printf("[Hook]post-deploy of service: \"%s\"\n", service.Name)
	} else {
		fmt.Printf("[Hook]post-deploy of global\n")
	}

	return h.handle(xHook.PostDeploy, createOptions, upOptions, pullOptions, service)
}

func (h *hook) handle(commands []types.ShellCommand, createOptions createOptions, upOptions upOptions, pullOptions pullOptions, service *types.ServiceConfig) error {
	for _, command := range commands {
		if exe := h.parseCommand(command, service); exe != nil {
			if err := exe.run(h, service); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *hook) parseCommand(command types.ShellCommand, service *types.ServiceConfig) *execute {
	if len(command) <= 0 {
		return nil
	}

	workDir := filepath.Dir(h.project.ComposeFiles[0]) // 相對docker-compose.yml文件的工作目錄
	var newCommand types.ShellCommand
	for _, arg := range command {
		if arg == "{ARGS}" {
			newCommand = append(newCommand, os.Args[2:]...)
		} else {
			newCommand = append(newCommand, arg)
		}
	}

	if len(newCommand) >= 2 {
		switch newCommand[0] {
		case "igo-key":
			path := filepath.Join(workDir, strings.TrimSpace(newCommand[1])+".gop")
			return &execute{
				path:        path,
				content:     h.getXKey(newCommand[1], service),
				executeType: igoKey,
				work:        workDir,
				command:     append(types.ShellCommand{"igop", path}, newCommand[2:]...),
			}
		case "igo-path":
			path := strings.TrimSpace(newCommand[1])
			if !filepath.IsAbs(path) { // 改為相對於docker-compose.yml文件的工作目錄
				path = filepath.Join(workDir, path)
			}
			return &execute{
				path:        path,
				content:     "",
				executeType: igoPath,
				work:        workDir,
				command:     append(types.ShellCommand{"igop", path}, newCommand[2:]...),
			}
		case "shell-key":
			path := filepath.Join(workDir, strings.TrimSpace(newCommand[1])+".sh")
			return &execute{
				path:        path,
				content:     h.getXKey(newCommand[1], service),
				executeType: shellKey,
				work:        workDir,
				command:     append(types.ShellCommand{"/usr/bin/sh", path}, newCommand[2:]...),
			}
		}
	}

	return &execute{
		executeType: shell,
		command:     command,
		work:        workDir,
	}
}

func (h *hook) getXKey(name string, service *types.ServiceConfig) string {
	name = strings.TrimSpace(name)
	if service != nil {
		if c, ok := service.Extensions[name]; ok {
			return c.(string)
		}
	}

	return h.project.Extensions[name].(string)
}

type executeType string

const (
	igoKey   executeType = "igo-key"
	igoPath              = "igo-path"
	shell                = "shell"
	shellKey             = "shell-key"
)

type execute struct {
	path        string
	content     string
	work        string
	executeType executeType
	command     types.ShellCommand
}

func (e *execute) run(h *hook, service *types.ServiceConfig) error {
	workDir := filepath.Dir(h.project.ComposeFiles[0]) // 相對docker-compose.yml文件的工作目錄

	fmt.Printf(" - execute %s %s: %+q\n", e.executeType, e.path, e.command)

	switch e.executeType {
	case igoKey:
		i := igo.IGo{
			Cmd:     h.cmd,
			Project: h.project,
			Service: service,
			Args:    e.command[2:],
		}
		return i.Run(e.path, e.content)
	case igoPath:
		i := igo.IGo{
			Cmd:     h.cmd,
			Project: h.project,
			Service: service,
			Args:    e.command[2:],
		}
		return i.RunPath(e.path)
	case shellKey:
		if err := ioutil.WriteFile(e.path, []byte(e.content), 0o644); err != nil {
			return err
		}
		fallthrough
	case shell:
		yamlBuf, _ := yaml.Marshal(h.project)
		var env []string
		for k, v := range h.project.Environment {
			env = append(env, k+"="+v)
		}

		cmd := exec.CommandContext(h.ctx, e.command[0], e.command[1:]...)
		cmd.Stdin = bytes.NewBuffer(yamlBuf)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = workDir
		cmd.Env = env
		return cmd.Run()
	}
	return nil
}
