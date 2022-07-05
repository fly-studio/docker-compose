// export by github.com/goplus/igop/cmd/qexp

//go:build go1.18
// +build go1.18

package nat

import (
	q "github.com/docker/go-connections/nat"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "nat",
		Path: "github.com/docker/go-connections/nat",
		Deps: map[string]string{
			"fmt":     "fmt",
			"net":     "net",
			"sort":    "sort",
			"strconv": "strconv",
			"strings": "strings",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]igop.NamedType{
			"Port":        {reflect.TypeOf((*q.Port)(nil)).Elem(), "Int,Port,Proto,Range", ""},
			"PortBinding": {reflect.TypeOf((*q.PortBinding)(nil)).Elem(), "", ""},
			"PortMap":     {reflect.TypeOf((*q.PortMap)(nil)).Elem(), "", ""},
			"PortMapping": {reflect.TypeOf((*q.PortMapping)(nil)).Elem(), "", ""},
			"PortSet":     {reflect.TypeOf((*q.PortSet)(nil)).Elem(), "", ""},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"NewPort":             reflect.ValueOf(q.NewPort),
			"ParsePort":           reflect.ValueOf(q.ParsePort),
			"ParsePortRange":      reflect.ValueOf(q.ParsePortRange),
			"ParsePortRangeToInt": reflect.ValueOf(q.ParsePortRangeToInt),
			"ParsePortSpec":       reflect.ValueOf(q.ParsePortSpec),
			"ParsePortSpecs":      reflect.ValueOf(q.ParsePortSpecs),
			"PartParser":          reflect.ValueOf(q.PartParser),
			"Sort":                reflect.ValueOf(q.Sort),
			"SortPortMap":         reflect.ValueOf(q.SortPortMap),
			"SplitProtoPort":      reflect.ValueOf(q.SplitProtoPort),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
