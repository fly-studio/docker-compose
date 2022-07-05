// export by github.com/goplus/igop/cmd/qexp

//go:build go1.18
// +build go1.18

package loader

import (
	q "github.com/compose-spec/compose-go/loader"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "loader",
		Path: "github.com/compose-spec/compose-go/loader",
		Deps: map[string]string{
			"fmt": "fmt",
			"github.com/compose-spec/compose-go/consts":        "consts",
			"github.com/compose-spec/compose-go/dotenv":        "dotenv",
			"github.com/compose-spec/compose-go/errdefs":       "errdefs",
			"github.com/compose-spec/compose-go/interpolation": "interpolation",
			"github.com/compose-spec/compose-go/schema":        "schema",
			"github.com/compose-spec/compose-go/template":      "template",
			"github.com/compose-spec/compose-go/types":         "types",
			"github.com/docker/go-units":                       "units",
			"github.com/imdario/mergo":                         "mergo",
			"github.com/mattn/go-shellwords":                   "shellwords",
			"github.com/mitchellh/mapstructure":                "mapstructure",
			"github.com/pkg/errors":                            "errors",
			"github.com/sirupsen/logrus":                       "logrus",
			"gopkg.in/yaml.v2":                                 "yaml",
			"io/ioutil":                                        "ioutil",
			"os":                                               "os",
			"path":                                             "path",
			"path/filepath":                                    "filepath",
			"reflect":                                          "reflect",
			"regexp":                                           "regexp",
			"sort":                                             "sort",
			"strconv":                                          "strconv",
			"strings":                                          "strings",
			"time":                                             "time",
			"unicode":                                          "unicode",
			"unicode/utf8":                                     "utf8",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]igop.NamedType{
			"ForbiddenPropertiesError": {reflect.TypeOf((*q.ForbiddenPropertiesError)(nil)).Elem(), "", "Error"},
			"Options":                  {reflect.TypeOf((*q.Options)(nil)).Elem(), "GetProjectName", "SetProjectName"},
			"Transformer":              {reflect.TypeOf((*q.Transformer)(nil)).Elem(), "", ""},
			"TransformerFunc":          {reflect.TypeOf((*q.TransformerFunc)(nil)).Elem(), "", ""},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"Propagations": reflect.ValueOf(&q.Propagations),
		},
		Funcs: map[string]reflect.Value{
			"Load":                 reflect.ValueOf(q.Load),
			"LoadConfigObjs":       reflect.ValueOf(q.LoadConfigObjs),
			"LoadNetworks":         reflect.ValueOf(q.LoadNetworks),
			"LoadSecrets":          reflect.ValueOf(q.LoadSecrets),
			"LoadService":          reflect.ValueOf(q.LoadService),
			"LoadServices":         reflect.ValueOf(q.LoadServices),
			"LoadVolumes":          reflect.ValueOf(q.LoadVolumes),
			"NormalizeProjectName": reflect.ValueOf(q.NormalizeProjectName),
			"ParseShortSSHSyntax":  reflect.ValueOf(q.ParseShortSSHSyntax),
			"ParseVolume":          reflect.ValueOf(q.ParseVolume),
			"ParseYAML":            reflect.ValueOf(q.ParseYAML),
			"Transform":            reflect.ValueOf(q.Transform),
			"WithDiscardEnvFiles":  reflect.ValueOf(q.WithDiscardEnvFiles),
			"WithSkipValidation":   reflect.ValueOf(q.WithSkipValidation),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
