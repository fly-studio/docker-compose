package compose

import (
	"context"
	"fmt"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/compose/v2/igo"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type xHooks struct {
	PreDeploy  types.ShellCommand `mapstructure:"pre-deploy"`
	PostDeploy types.ShellCommand `mapstructure:"post-deploy"`
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

func (h *hook) parseIGo(command string) (key string, content string, err error) {
	if strings.HasPrefix(command, "igo-key:") {
		if c, ok := h.project.Extensions[command[8:]].(string); ok {
			return command[8:], c, nil
		}
		return "", "", fmt.Errorf("igo-key \"%s\" is not exists, or invalid string. key must starts with \"x-\"\n", command[4:])
	} else if strings.HasPrefix(command, "igo-path:") {
		return command[9:], "", nil
	}

	return "", "", nil
}

func (h *hook) executeIGo(vpath string, content string, createOptions createOptions, upOptions upOptions, pullOptions pullOptions, services []string) error {
	i := igo.IGo{
		Cmd:      h.cmd,
		Project:  h.project,
		Services: services,
	}
	return i.Run(vpath, content)
}

func (h *hook) executeIGoPath(path string, createOptions createOptions, upOptions upOptions, pullOptions pullOptions, services []string) error {
	i := igo.IGo{
		Cmd:      h.cmd,
		Project:  h.project,
		Services: services,
	}

	return i.RunPath(path)
}

func (h *hook) PreDeploy(createOptions createOptions, upOptions upOptions, pullOptions pullOptions, services []string) error {
	return h.handle(h.xHooks.PreDeploy, createOptions, upOptions, pullOptions, services)
}

func (h *hook) PostDeploy(createOptions createOptions, upOptions upOptions, pullOptions pullOptions, services []string) error {
	return h.handle(h.xHooks.PostDeploy, createOptions, upOptions, pullOptions, services)
}

func (h *hook) handle(command types.ShellCommand, createOptions createOptions, upOptions upOptions, pullOptions pullOptions, services []string) error {
	if len(command) <= 0 {
		return nil
	}

	workDir := filepath.Dir(h.project.ComposeFiles[0]) // 相對docker-compose.yml文件的工作目錄
	// 檢查是否是igo
	if len(command) == 1 && strings.HasPrefix(command[0], "igo-") {
		key, content, err := h.parseIGo(command[0])
		if err != nil {
			return err
		}

		if content != "" {
			return h.executeIGo(filepath.Join(workDir, key+".gop"), content, createOptions, upOptions, pullOptions, services)
		} else if key != "" {
			if !filepath.IsAbs(key) { // 改為相對於docker-compose.yml文件的工作目錄
				key = filepath.Join(workDir, key)
			}
			return h.executeIGoPath(key, createOptions, upOptions, pullOptions, services)
		}
	}

	cmd := exec.CommandContext(h.ctx, command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = workDir
	return cmd.Run()
}
