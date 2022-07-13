package igo

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/goplus/igop"
	"github.com/spf13/cobra"
	"reflect"
)

import _ "github.com/goplus/igop/pkg"
import _ "github.com/docker/compose/v2/igo/pkgs"

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "igo",
		Path: "igo",
		Deps: map[string]string{
			"github.com/compose-spec/compose-go/types": "types",
			"github.com/spf13/cobra":                   "cobra",
			"github.com/bitly/go-simplejson":           "simplejson",
		},
		Interfaces: map[string]reflect.Type{},
		AliasTypes: map[string]reflect.Type{},
		NamedTypes: map[string]igop.NamedType{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"GetCmd":     reflect.ValueOf(GetCmd),
			"GetService": reflect.ValueOf(GetService),
			"GetProject": reflect.ValueOf(GetProject),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}

type IGo struct {
	Cmd     *cobra.Command
	Project *types.Project
	Service *types.ServiceConfig
	Args    types.ShellCommand
}

var globalIGo IGo

func GetCmd() *cobra.Command {
	return globalIGo.Cmd
}

func GetService() *types.ServiceConfig {
	return globalIGo.Service
}

func GetProject() *types.Project {
	return globalIGo.Project
}

func (i *IGo) Run(vpath string, content string) error {
	// 暫時沒有處理多線程下的運行衝突問題
	globalIGo = *i

	if vpath == "" {
		vpath = "main.gop"
	}
	_, err := igop.RunFile(vpath, content, i.Args, 0)
	return err
}

func (i *IGo) RunPath(path string) error {
	// 暫時沒有處理多線程下的運行衝突問題
	globalIGo = *i

	_, err := igop.Run(path, i.Args, 0)
	return err
}
