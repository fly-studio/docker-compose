// export by github.com/goplus/igop/cmd/qexp

//go:build go1.18
// +build go1.18

package mapstructure

import (
	q "github.com/mitchellh/mapstructure"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "mapstructure",
		Path: "github.com/mitchellh/mapstructure",
		Deps: map[string]string{
			"encoding":      "encoding",
			"encoding/json": "json",
			"errors":        "errors",
			"fmt":           "fmt",
			"net":           "net",
			"reflect":       "reflect",
			"sort":          "sort",
			"strconv":       "strconv",
			"strings":       "strings",
			"time":          "time",
		},
		Interfaces: map[string]reflect.Type{
			"DecodeHookFunc": reflect.TypeOf((*q.DecodeHookFunc)(nil)).Elem(),
		},
		NamedTypes: map[string]igop.NamedType{
			"DecodeHookFuncKind":  {reflect.TypeOf((*q.DecodeHookFuncKind)(nil)).Elem(), "", ""},
			"DecodeHookFuncType":  {reflect.TypeOf((*q.DecodeHookFuncType)(nil)).Elem(), "", ""},
			"DecodeHookFuncValue": {reflect.TypeOf((*q.DecodeHookFuncValue)(nil)).Elem(), "", ""},
			"Decoder":             {reflect.TypeOf((*q.Decoder)(nil)).Elem(), "", "Decode,decode,decodeArray,decodeBasic,decodeBool,decodeFloat,decodeFunc,decodeInt,decodeMap,decodeMapFromMap,decodeMapFromSlice,decodeMapFromStruct,decodePtr,decodeSlice,decodeString,decodeStruct,decodeStructFromMap,decodeUint"},
			"DecoderConfig":       {reflect.TypeOf((*q.DecoderConfig)(nil)).Elem(), "", ""},
			"Error":               {reflect.TypeOf((*q.Error)(nil)).Elem(), "", "Error,WrappedErrors"},
			"Metadata":            {reflect.TypeOf((*q.Metadata)(nil)).Elem(), "", ""},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars:       map[string]reflect.Value{},
		Funcs: map[string]reflect.Value{
			"ComposeDecodeHookFunc":        reflect.ValueOf(q.ComposeDecodeHookFunc),
			"Decode":                       reflect.ValueOf(q.Decode),
			"DecodeHookExec":               reflect.ValueOf(q.DecodeHookExec),
			"DecodeMetadata":               reflect.ValueOf(q.DecodeMetadata),
			"NewDecoder":                   reflect.ValueOf(q.NewDecoder),
			"OrComposeDecodeHookFunc":      reflect.ValueOf(q.OrComposeDecodeHookFunc),
			"RecursiveStructToMapHookFunc": reflect.ValueOf(q.RecursiveStructToMapHookFunc),
			"StringToIPHookFunc":           reflect.ValueOf(q.StringToIPHookFunc),
			"StringToIPNetHookFunc":        reflect.ValueOf(q.StringToIPNetHookFunc),
			"StringToSliceHookFunc":        reflect.ValueOf(q.StringToSliceHookFunc),
			"StringToTimeDurationHookFunc": reflect.ValueOf(q.StringToTimeDurationHookFunc),
			"StringToTimeHookFunc":         reflect.ValueOf(q.StringToTimeHookFunc),
			"TextUnmarshallerHookFunc":     reflect.ValueOf(q.TextUnmarshallerHookFunc),
			"WeakDecode":                   reflect.ValueOf(q.WeakDecode),
			"WeakDecodeMetadata":           reflect.ValueOf(q.WeakDecodeMetadata),
			"WeaklyTypedHook":              reflect.ValueOf(q.WeaklyTypedHook),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
