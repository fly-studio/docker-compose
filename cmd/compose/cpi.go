package compose

import (
	"context"
	"fmt"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/compose/v2/pkg/api"
	cliTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/system"
	"github.com/spf13/cobra"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func cpiCommand(p *projectOptions, dockerCli command.Cli, backend api.Service) *cobra.Command {
	opts := copyOptions{
		projectOptions: p,
	}
	cpiCmd := &cobra.Command{
		Use:   "cpi [OPTIONS] [SERVICE] [PATH_IN_IMAGE] [LOCAL_PATH]",
		Short: "Copy path from image of service to local filesystem",
		Args:  cli.ExactArgs(3),
		PreRunE: AdaptCmd(func(ctx context.Context, cmd *cobra.Command, args []string) error {
			return nil
		}),
		RunE: Adapt(func(ctx context.Context, args []string) error {
			service := args[0]
			source := args[1]
			destination := args[2]

			project, err := p.toProject([]string{service})
			if err != nil {
				return err
			} else if len(project.Services) != 1 {
				return fmt.Errorf("service not found: %s", service)
			}

			return runCpi(ctx, dockerCli, project, opts, source, destination)
		}),
		ValidArgsFunction: serviceCompletion(p),
	}

	flags := cpiCmd.Flags()
	flags.BoolVarP(&opts.followLink, "follow-link", "L", false, "Always follow symbol link in PATH_IN_IMAGE")

	return cpiCmd
}

func runCpi(ctx context.Context, dockerCli command.Cli, project *types.Project, opts copyOptions, source, destination string) error {
	rand.Seed(time.Now().UnixNano())
	service := project.Services[0]
	name := service.Name + "-cp-temp-" + strconv.Itoa(rand.Intn(1000000))

	// https://stackoverflow.com/questions/25292198/docker-how-can-i-copy-a-file-from-an-image-to-a-host
	fmt.Printf(" - Create temporary container \"%s\" of image: [%s]\n", name, service.Image)
	if _, err := dockerCli.Client().ContainerCreate(ctx, &container.Config{
		Env:             []string{},
		Cmd:             strslice.StrSlice(service.Command),
		Image:           service.Image,
		WorkingDir:      service.WorkingDir,
		Entrypoint:      strslice.StrSlice(service.Entrypoint),
		NetworkDisabled: true,
		Labels:          service.Labels,
	}, nil, nil, nil, name); err != nil {
		return err
	}

	fmt.Printf(" - Copy %s of image: [%s] to %s\n", source, service.Image, destination)
	if err := os.MkdirAll(destination, 0644); err != nil {
		return err
	}
	if err := copyFromContainer(ctx, dockerCli, name, source, destination, opts.followLink); err != nil {
		return err
	}

	fmt.Printf(" - Remove temporary container \"%s\" of image: [%s]\n", name, service.Image)
	if err := dockerCli.Client().ContainerRemove(ctx, name, cliTypes.ContainerRemoveOptions{}); err != nil {
		return err
	}

	fmt.Printf("Copy successful from image: [%s]", service.Image)

	return nil
}

func copyFromContainer(ctx context.Context, dockerCli command.Cli, containerID, srcPath, dstPath string, followLink bool) error {
	var err error
	if dstPath != "-" {
		// Get an absolute destination path.
		dstPath, err = resolveLocalPath(dstPath)
		if err != nil {
			return err
		}
	}

	if err := command.ValidateOutputPath(dstPath); err != nil {
		return err
	}

	// if client requests to follow symbol link, then must decide target file to be copied
	var rebaseName string
	if followLink {
		srcStat, err := dockerCli.Client().ContainerStatPath(ctx, containerID, srcPath)

		// If the destination is a symbolic link, we should follow it.
		if err == nil && srcStat.Mode&os.ModeSymlink != 0 {
			linkTarget := srcStat.LinkTarget
			if !system.IsAbs(linkTarget) {
				// Join with the parent directory.
				srcParent, _ := archive.SplitPathDirEntry(srcPath)
				linkTarget = filepath.Join(srcParent, linkTarget)
			}

			linkTarget, rebaseName = archive.GetRebaseName(srcPath, linkTarget)
			srcPath = linkTarget
		}
	}

	content, stat, err := dockerCli.Client().CopyFromContainer(ctx, containerID, srcPath)
	if err != nil {
		return err
	}
	defer content.Close() //nolint:errcheck

	if dstPath == "-" {
		_, err = io.Copy(os.Stdout, content)
		return err
	}

	srcInfo := archive.CopyInfo{
		Path:       srcPath,
		Exists:     true,
		IsDir:      stat.Mode.IsDir(),
		RebaseName: rebaseName,
	}

	preArchive := content
	if len(srcInfo.RebaseName) != 0 {
		_, srcBase := archive.SplitPathDirEntry(srcInfo.Path)
		preArchive = archive.RebaseArchiveEntries(content, srcBase, srcInfo.RebaseName)
	}

	return archive.CopyTo(preArchive, srcInfo, dstPath)
}

func resolveLocalPath(localPath string) (absPath string, err error) {
	if absPath, err = filepath.Abs(localPath); err != nil {
		return
	}
	return archive.PreserveTrailingDotOrSeparator(absPath, localPath, filepath.Separator), nil
}
