// export by github.com/goplus/igop/cmd/qexp

//go:build go1.18
// +build go1.18

package simplejson

import (
	q "github.com/bitly/go-simplejson"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "simplejson",
		Path: "github.com/bitly/go-simplejson",
		Deps: map[string]string{
			"bytes":         "bytes",
			"encoding/json": "json",
			"errors":        "errors",
			"io":            "io",
			"log":           "log",
			"reflect":       "reflect",
			"strconv":       "strconv",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]igop.NamedType{
			"Json": {reflect.TypeOf((*q.Json)(nil)).Elem(), "", "Array,Bool,Bytes,CheckGet,Del,Encode,EncodePretty,Float64,Get,GetIndex,GetPath,Int,Int64,Interface,Map,MarshalJSON,MustArray,MustBool,MustFloat64,MustInt,MustInt64,MustMap,MustString,MustStringArray,MustUint64,Set,SetPath,String,StringArray,Uint64,UnmarshalJSON"},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"New":           reflect.ValueOf(q.New),
			"NewFromReader": reflect.ValueOf(q.NewFromReader),
			"NewJson":       reflect.ValueOf(q.NewJson),
			"Version":       reflect.ValueOf(q.Version),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
