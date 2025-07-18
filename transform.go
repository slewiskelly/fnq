package fnq

import (
	"context"
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/encoding/yaml"
	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"

	"github.com/slewiskelly/fnq/internal/pkg/mod"
)

// Transform returns a ResourceListProcessorFunc for transforming resources
// within a resource list.
//
// Each resource within the resource list will be transformed according to a
// corresponding resource transformer in the given module.
//
// See the [CUE reference] page for how validators should be expressed.
//
// A panic will occur if there is an error while retrieving the module.
//
// [CUE reference]: https://github.com/slewiskelly/fnq/docs/references/cue.md
func Transform(ctx context.Context, module string) fn.ResourceListProcessorFunc {
	v, err := mod.Get(ctx, module)
	if err != nil {
		panic(err)
	}

	return (&processor{v: v}).transform
}

func (p *processor) transform(rl *fn.ResourceList) (bool, error) {
	var results fn.Results

	for _, obj := range rl.Items.WhereNot(fn.HasAnnotations(map[string]string{"no-transform": "true"})) {
		i, ok, err := p.xform(obj)
		if err != nil {
			if errors.Is(err, errNotFound) {
				continue
			}

			results = append(results, fn.ErrorConfigObjectResult(err, obj))
			continue
		}

		if ok {
			rl.Results.Infof("Transformed %s", obj.GetId())
		} else {
			rl.Results.Infof("Unchanged %s", obj.GetId())
		}

		rl.UpsertObjectToItems(i, nil, true)
	}

	rl.Results = append(rl.Results, results...)

	if len(results) > 0 {
		return false, errors.New("transformation(s) failed")
	}

	return true, nil
}

func (p *processor) xform(obj *fn.KubeObject) (*fn.KubeObject, bool, error) {
	grp, ver, kind := gvk(obj)

	v := p.v.LookupPath(cue.ParsePath(fmt.Sprintf("Transformers[%q][%q][%q]", grp, ver, kind)))
	if err := v.Err(); err != nil {
		if !v.Exists() {
			return nil, false, errNotFound
		}

		return nil, false, errDetails(err)
	}

	a, err := yaml.Extract("", obj.String())
	if err != nil {
		return nil, false, errDetails(err)
	}

	w := v.Context().BuildFile(a)
	if err := w.Err(); err != nil {
		return nil, false, errDetails(err)
	}

	v = v.Unify(w)
	if err := v.Validate(cue.Final(), cue.Hidden(false)); err != nil {
		return nil, false, errDetails(err)
	}

	b, err := yaml.Encode(v)
	if err != nil {
		return nil, false, errDetails(err)
	}

	o, err := fn.ParseKubeObject(b)

	return o, obj.String() != string(b), err
}
