package compose

import (
	"bytes"
	"context"
	"fmt"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
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
	ctx     context.Context
	cmd     *cobra.Command
	project *types.Project
	backend api.Service
	xHooks  xHooks
}

func (h *hook) parse() error {
	return loader.Transform(h.project.Extensions["x-hooks"], &h.xHooks)
}

func (h *hook) PreDeploy(createOptions createOptions, upOptions upOptions, pullOptions pullOptions, services []string) error {
	return h.handle(h.xHooks.PreDeploy, createOptions, upOptions, pullOptions, services)
}

func (h *hook) PostDeploy(createOptions createOptions, upOptions upOptions, pullOptions pullOptions, services []string) error {
	return h.handle(h.xHooks.PostDeploy, createOptions, upOptions, pullOptions, services)
}

func (h *hook) handle(commands []types.ShellCommand, createOptions createOptions, upOptions upOptions, pullOptions pullOptions, services []string) error {
	for _, command := range commands {
		if exe := h.parseCommand(command); exe != nil {

			if err := exe.run(h, services); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *hook) parseCommand(command types.ShellCommand) *execute {
	if len(command) <= 0 {
		return nil
	}

	workDir := filepath.Dir(h.project.ComposeFiles[0]) // 相對docker-compose.yml文件的工作目錄

	if len(command) == 2 {
		command[1] = strings.TrimSpace(command[1])
		switch command[0] {
		case "igo-key":
			return &execute{
				path:        filepath.Join(workDir, command[1]+".gop"),
				content:     h.project.Extensions[command[1]].(string),
				executeType: igoKey,
				work:        workDir,
			}
		case "igo-path":
			path := command[1]
			if !filepath.IsAbs(path) { // 改為相對於docker-compose.yml文件的工作目錄
				path = filepath.Join(workDir, path)
			}
			return &execute{
				path:        path,
				content:     "",
				executeType: igoPath,
				work:        workDir,
			}
		case "shell-key":
			path := filepath.Join(workDir, command[1]+".sh")
			return &execute{
				path:        path,
				content:     h.project.Extensions[command[1]].(string),
				executeType: shellKey,
				command:     append(types.ShellCommand{"/usr/bin/sh", path}, os.Args[2:]...),
				work:        workDir,
			}
		}
	}

	return &execute{
		executeType: shell,
		command:     command,
		work:        workDir,
	}
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

func (e *execute) run(h *hook, services []string) error {
	workDir := filepath.Dir(h.project.ComposeFiles[0]) // 相對docker-compose.yml文件的工作目錄

	fmt.Printf("execute %s %s: %+q\n", e.executeType, e.path, e.command)

	switch e.executeType {
	case igoKey:
		i := igo.IGo{
			Cmd:      h.cmd,
			Project:  h.project,
			Services: services,
		}
		return i.Run(e.path, e.content)
	case igoPath:
		i := igo.IGo{
			Cmd:      h.cmd,
			Project:  h.project,
			Services: services,
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
