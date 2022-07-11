//go:build go1.18
// +build go1.18

package pkgs

/*
qexp -outdir . -addtags "//+build go1.18" -filename go118_export \
github.com/spf13/pflag \
github.com/spf13/cobra \
github.com/compose-spec/compose-go/types \
github.com/compose-spec/compose-go/loader \
golang.org/x/sync/errgroup \
github.com/opencontainers/go-digest \
github.com/mitchellh/mapstructure \
github.com/docker/go-connections/nat \
github.com/distribution/distribution/v3/reference \
github.com/distribution/distribution/v3/digestset \
github.com/bitly/go-simplejson
*/

import (
	_ "github.com/docker/compose/v2/igo/pkgs/github.com/bitly/go-simplejson"
	_ "github.com/docker/compose/v2/igo/pkgs/github.com/compose-spec/compose-go/types"
	_ "github.com/docker/compose/v2/igo/pkgs/github.com/distribution/distribution/v3/digestset"
	_ "github.com/docker/compose/v2/igo/pkgs/github.com/distribution/distribution/v3/reference"
	_ "github.com/docker/compose/v2/igo/pkgs/github.com/docker/go-connections/nat"
	_ "github.com/docker/compose/v2/igo/pkgs/github.com/mitchellh/mapstructure"
	_ "github.com/docker/compose/v2/igo/pkgs/github.com/opencontainers/go-digest"
	_ "github.com/docker/compose/v2/igo/pkgs/github.com/spf13/cobra"
	_ "github.com/docker/compose/v2/igo/pkgs/github.com/spf13/pflag"
	_ "github.com/docker/compose/v2/igo/pkgs/golang.org/x/sync/errgroup"
)
