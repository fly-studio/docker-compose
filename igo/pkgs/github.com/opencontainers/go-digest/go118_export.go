// export by github.com/goplus/igop/cmd/qexp

//go:build go1.18
// +build go1.18

package digest

import (
	q "github.com/opencontainers/go-digest"

	"go/constant"
	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "digest",
		Path: "github.com/opencontainers/go-digest",
		Deps: map[string]string{
			"crypto":  "crypto",
			"fmt":     "fmt",
			"hash":    "hash",
			"io":      "io",
			"regexp":  "regexp",
			"strings": "strings",
		},
		Interfaces: map[string]reflect.Type{
			"Digester": reflect.TypeOf((*q.Digester)(nil)).Elem(),
			"Verifier": reflect.TypeOf((*q.Verifier)(nil)).Elem(),
		},
		NamedTypes: map[string]igop.NamedType{
			"Algorithm": {reflect.TypeOf((*q.Algorithm)(nil)).Elem(), "Available,Digester,Encode,FromBytes,FromReader,FromString,Hash,Size,String,Validate", "Set"},
			"Digest":    {reflect.TypeOf((*q.Digest)(nil)).Elem(), "Algorithm,Encoded,Hex,String,Validate,Verifier,sepIndex", ""},
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"DigestRegexp":           reflect.ValueOf(&q.DigestRegexp),
			"DigestRegexpAnchored":   reflect.ValueOf(&q.DigestRegexpAnchored),
			"ErrDigestInvalidFormat": reflect.ValueOf(&q.ErrDigestInvalidFormat),
			"ErrDigestInvalidLength": reflect.ValueOf(&q.ErrDigestInvalidLength),
			"ErrDigestUnsupported":   reflect.ValueOf(&q.ErrDigestUnsupported),
		},
		Funcs: map[string]reflect.Value{
			"FromBytes":            reflect.ValueOf(q.FromBytes),
			"FromReader":           reflect.ValueOf(q.FromReader),
			"FromString":           reflect.ValueOf(q.FromString),
			"NewDigest":            reflect.ValueOf(q.NewDigest),
			"NewDigestFromBytes":   reflect.ValueOf(q.NewDigestFromBytes),
			"NewDigestFromEncoded": reflect.ValueOf(q.NewDigestFromEncoded),
			"NewDigestFromHex":     reflect.ValueOf(q.NewDigestFromHex),
			"Parse":                reflect.ValueOf(q.Parse),
		},
		TypedConsts: map[string]igop.TypedConst{
			"Canonical": {reflect.TypeOf(q.Canonical), constant.MakeString(string(q.Canonical))},
			"SHA256":    {reflect.TypeOf(q.SHA256), constant.MakeString(string(q.SHA256))},
			"SHA384":    {reflect.TypeOf(q.SHA384), constant.MakeString(string(q.SHA384))},
			"SHA512":    {reflect.TypeOf(q.SHA512), constant.MakeString(string(q.SHA512))},
		},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
