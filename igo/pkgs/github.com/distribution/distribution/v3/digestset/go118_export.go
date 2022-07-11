// export by github.com/goplus/igop/cmd/qexp

//go:build go1.18
// +build go1.18

package digestset

import (
	q "github.com/distribution/distribution/v3/digestset"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "digestset",
		Path: "github.com/distribution/distribution/v3/digestset",
		Deps: map[string]string{
			"errors":                              "errors",
			"github.com/opencontainers/go-digest": "digest",
			"sort":                                "sort",
			"strings":                             "strings",
			"sync":                                "sync",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]igop.NamedType{
			"Set": {reflect.TypeOf((*q.Set)(nil)).Elem(), "", "Add,All,Lookup,Remove"},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"ErrDigestAmbiguous": reflect.ValueOf(&q.ErrDigestAmbiguous),
			"ErrDigestNotFound":  reflect.ValueOf(&q.ErrDigestNotFound),
		},
		Funcs: map[string]reflect.Value{
			"NewSet":         reflect.ValueOf(q.NewSet),
			"ShortCodeTable": reflect.ValueOf(q.ShortCodeTable),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
