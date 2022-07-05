// export by github.com/goplus/igop/cmd/qexp

//go:build go1.18
// +build go1.18

package reference

import (
	q "github.com/distribution/distribution/v3/reference"

	"go/constant"
	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "reference",
		Path: "github.com/distribution/distribution/v3/reference",
		Deps: map[string]string{
			"errors": "errors",
			"fmt":    "fmt",
			"github.com/distribution/distribution/v3/digestset": "digestset",
			"github.com/opencontainers/go-digest":               "digest",
			"path":                                              "path",
			"regexp":                                            "regexp",
			"strings":                                           "strings",
		},
		Interfaces: map[string]reflect.Type{
			"Canonical":   reflect.TypeOf((*q.Canonical)(nil)).Elem(),
			"Digested":    reflect.TypeOf((*q.Digested)(nil)).Elem(),
			"Named":       reflect.TypeOf((*q.Named)(nil)).Elem(),
			"NamedTagged": reflect.TypeOf((*q.NamedTagged)(nil)).Elem(),
			"Reference":   reflect.TypeOf((*q.Reference)(nil)).Elem(),
			"Tagged":      reflect.TypeOf((*q.Tagged)(nil)).Elem(),
		},
		NamedTypes: map[string]igop.NamedType{
			"Field": {reflect.TypeOf((*q.Field)(nil)).Elem(), "MarshalText,Reference", "UnmarshalText"},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"DigestRegexp":              reflect.ValueOf(&q.DigestRegexp),
			"DomainRegexp":              reflect.ValueOf(&q.DomainRegexp),
			"ErrDigestInvalidFormat":    reflect.ValueOf(&q.ErrDigestInvalidFormat),
			"ErrNameContainsUppercase":  reflect.ValueOf(&q.ErrNameContainsUppercase),
			"ErrNameEmpty":              reflect.ValueOf(&q.ErrNameEmpty),
			"ErrNameNotCanonical":       reflect.ValueOf(&q.ErrNameNotCanonical),
			"ErrNameTooLong":            reflect.ValueOf(&q.ErrNameTooLong),
			"ErrReferenceInvalidFormat": reflect.ValueOf(&q.ErrReferenceInvalidFormat),
			"ErrTagInvalidFormat":       reflect.ValueOf(&q.ErrTagInvalidFormat),
			"IdentifierRegexp":          reflect.ValueOf(&q.IdentifierRegexp),
			"NameRegexp":                reflect.ValueOf(&q.NameRegexp),
			"ReferenceRegexp":           reflect.ValueOf(&q.ReferenceRegexp),
			"ShortIdentifierRegexp":     reflect.ValueOf(&q.ShortIdentifierRegexp),
			"TagRegexp":                 reflect.ValueOf(&q.TagRegexp),
		},
		Funcs: map[string]reflect.Value{
			"AsField":                  reflect.ValueOf(q.AsField),
			"Domain":                   reflect.ValueOf(q.Domain),
			"FamiliarMatch":            reflect.ValueOf(q.FamiliarMatch),
			"FamiliarName":             reflect.ValueOf(q.FamiliarName),
			"FamiliarString":           reflect.ValueOf(q.FamiliarString),
			"IsNameOnly":               reflect.ValueOf(q.IsNameOnly),
			"Parse":                    reflect.ValueOf(q.Parse),
			"ParseAnyReference":        reflect.ValueOf(q.ParseAnyReference),
			"ParseAnyReferenceWithSet": reflect.ValueOf(q.ParseAnyReferenceWithSet),
			"ParseDockerRef":           reflect.ValueOf(q.ParseDockerRef),
			"ParseNamed":               reflect.ValueOf(q.ParseNamed),
			"ParseNormalizedNamed":     reflect.ValueOf(q.ParseNormalizedNamed),
			"Path":                     reflect.ValueOf(q.Path),
			"SplitHostname":            reflect.ValueOf(q.SplitHostname),
			"TagNameOnly":              reflect.ValueOf(q.TagNameOnly),
			"TrimNamed":                reflect.ValueOf(q.TrimNamed),
			"WithDigest":               reflect.ValueOf(q.WithDigest),
			"WithName":                 reflect.ValueOf(q.WithName),
			"WithTag":                  reflect.ValueOf(q.WithTag),
		},
		TypedConsts: map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{
			"NameTotalLengthMax": {"untyped int", constant.MakeInt64(int64(q.NameTotalLengthMax))},
		},
	})
}
