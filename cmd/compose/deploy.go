/*
   Copyright 2020 Docker Compose CLI authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package compose

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/compose-spec/compose-go/types"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/utils"
)

type deployOptions struct {
	upOptions

	serverId int
	force    bool
}

func deployCommand(p *projectOptions, backend api.Service) *cobra.Command {
	deploy := deployOptions{}
	create := createOptions{}
	deployCmd := &cobra.Command{
		Use:   "deploy [SERVICE...]",
		Short: "Deploy containers(make dirs, download configs, replace server data, and start containers)",
		PreRunE: AdaptCmd(func(ctx context.Context, cmd *cobra.Command, args []string) error {
			create.timeChanged = cmd.Flags().Changed("timeout")
			return validateFlags(&deploy.upOptions, &create)
		}),
		RunE: p.WithServices(func(ctx context.Context, project *types.Project, services []string) error {
			create.ignoreOrphans = utils.StringToBool(project.Environment["COMPOSE_IGNORE_ORPHANS"])
			if create.ignoreOrphans && create.removeOrphans {
				return fmt.Errorf("COMPOSE_IGNORE_ORPHANS and --remove-orphans cannot be combined")
			}
			return runDeploy(ctx, backend, create, deploy, project, services)
		}),
		ValidArgsFunction: serviceCompletion(p),
	}
	flags := deployCmd.Flags()
	flags.AddFlagSet(upCommand(p, backend).Flags())

	flags.IntVarP(&deploy.serverId, "server-id", "s", 0, "the serverId of current host os.")
	flags.BoolVarP(&deploy.force, "force", "f", false, "force download config.")

	return deployCmd
}

func runDeploy(ctx context.Context, backend api.Service, createOptions createOptions, deployOptions deployOptions, project *types.Project, services []string) error {
	if len(project.Services) == 0 {
		return fmt.Errorf("no service selected")
	}

	var serverApi string
	if serverApi, _ = project.Environment["SERVER_API"]; serverApi == "" {
		return fmt.Errorf("miss \"SERVER_API\" in Environment")
	} else if _, err := url.ParseRequestURI(serverApi); err != nil {
		return fmt.Errorf("invalid \"SERVER_API\" in Environment: %w", err)
	}

	for _, service := range project.Services {
		fmt.Println("initial services:", service.Name)
		for _, conf := range service.Configs {
			fmt.Println("\t", "read config:", conf.Source)
			if err := downloadConfig(&deployOptions, project, conf.Source); err != nil {
				return fmt.Errorf("read config \"%s\" error: %w", conf.Source, err)
			}
		}
	}

	return runUp(ctx, backend, createOptions, deployOptions.upOptions, project, services)
}

func downloadConfig(deployOptions *deployOptions, project *types.Project, name string) error {
	conf, ok := project.Configs[name]
	if !ok {
		return fmt.Errorf("config \"%s\" not exists", name)
	}

	// 存在則不下載，除非強制下載
	if !deployOptions.force && pathExists(conf.File) {
		return nil
	}

	// 需要下載配置
	if xDownloadMap, ok := conf.Extensions["x-download"]; ok {
		var xDownload xDownload
		if err := mapstructure.Decode(xDownloadMap, &xDownload); err != nil {
			return err
		}
		fmt.Println("\t", "download config file:", conf.File)
		if err := download(xDownload, conf.File); err != nil {
			return err
		}

	}
	// 沒有傳遞 --server-id
	if deployOptions.serverId <= 0 {
		return nil
	}
	// 需要替換server的內容，只支援json格式
	if _, ok = conf.Extensions["x-standalone-server"]; ok {
		fmt.Println("\t", "get/replace server information to config file:", conf.File, "of server-id:", deployOptions.serverId)
		return replaceServerInformation(conf.File, project.Environment["SERVER_API"], deployOptions.serverId)
	}

	return nil
}

type xDownload struct {
	Url     string            `mapstructure:"url"`
	Headers map[string]string `mapstructure:"headers"`
}

func download(xDownload xDownload, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o644); err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodGet, xDownload.Url, nil)
	if err != nil {
		return err
	}
	for k, v := range xDownload.Headers {
		request.Header.Set(k, v)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, response.Body)
	return err
}

func replaceServerInformation(jsonFile, serverApi string, serverId int) error {
	// 遠端讀取server信息
	response, err := http.Get(strings.ReplaceAll(serverApi, "{server_id}", strconv.Itoa(serverId)))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	server, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	} else if !bytes.Contains(server, []byte("\"name\"")) {
		return fmt.Errorf("get server information error: %d", serverId)
	}
	// 讀取配置文件並替換
	f, err := os.OpenFile(jsonFile, os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	// 讀取配置
	buf, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	buf = bytes.ReplaceAll(buf, []byte("{-server-}"), server)

	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if err = f.Truncate(0); err != nil {
		return err
	}
	if _, err = f.Write(buf); err != nil {
		return err
	}

	return nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return err == nil
}
