// Package mod provides provides functionality for retrieving CUE modules.
package mod

import (
	"context"
	"os"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/mod/modconfig"
	"cuelang.org/go/mod/module"
)

// Get retrieves and loads a CUE module located at the specified path.
//
// The path must be in the form of $MODULE@$VERSION, where $VERSION is canonical.
func Get(ctx context.Context, path string) (cue.Value, error) {
	ver, err := module.ParseVersion(path)
	if err != nil {
		return cue.Value{}, err
	}

	reg, err := modconfig.NewRegistry(nil)
	if err != nil {
		return cue.Value{}, err
	}

	mod, err := reg.Fetch(ctx, ver)
	if err != nil {
		return cue.Value{}, err
	}

	tmp := filepath.Join(os.TempDir(), ver.String())

	if err := os.MkdirAll(tmp, 0700); err != nil {
		return cue.Value{}, err
	}
	defer os.RemoveAll(tmp)

	os.CopyFS(tmp, mod.FS)

	insts := load.Instances([]string{"."}, &load.Config{
		Dir:      tmp,
		Package:  "_",
		Registry: reg,
	})

	v := cuecontext.New().BuildInstance(insts[0])
	if err := v.Err(); err != nil {
		return cue.Value{}, err
	}

	return v, nil
}
