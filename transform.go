package fnq

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/encoding/yaml"
	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
)

// Transform returns a ResourceListProcessorFunc for transforming resources
// within a resource list.
//
// Each resource within the resource list will be transformed according to a
// corresponding resource transformer.
//
// See the [CUE reference] page for how validators should be expressed.
//
// [CUE reference]: https://github.com/slewiskelly/fnq/docs/references/cue.md
func Transform(v cue.Value) fn.ResourceListProcessorFunc {
	if err := v.Err(); err != nil {
		panic(fmt.Errorf("invalid transformer: %w", err))
	}

	return (&processor{v: v}).transform
}

func (p *processor) transform(rl *fn.ResourceList) (bool, error) {
	var results fn.Results

	for _, obj := range rl.Items.WhereNot(fn.HasAnnotations(map[string]string{"no-transform": "true"})) {
		i, err := p.xform(obj)
		if err != nil {
			if errors.Is(err, errNotFound) {
				continue
			}

			results = append(results, fn.ErrorConfigObjectResult(err, obj))
			continue
		}

		rl.UpsertObjectToItems(i, nil, true)
	}

	rl.Results = append(rl.Results, results...)

	if len(results) > 0 {
		return false, errors.New("transformation(s) failed")
	}

	return true, nil
}

func (p *processor) xform(obj *fn.KubeObject) (*fn.KubeObject, error) {
	v := p.v.LookupPath(cue.ParsePath("Transformers")).LookupPath(cue.ParsePath(gvk(obj)))
	if err := v.Err(); err != nil {
		if !v.Exists() {
			return nil, errNotFound
		}

		return nil, errDetails(err)
	}

	a, err := yaml.Extract("", obj.String())
	if err != nil {
		return nil, errDetails(err)
	}

	w := v.Context().BuildFile(a)
	if err := w.Err(); err != nil {
		return nil, errDetails(err)
	}

	v = v.Unify(w)
	if err := v.Validate(cue.Final(), cue.Hidden(false)); err != nil {
		return nil, errDetails(err)
	}

	b, err := yaml.Encode(v)
	if err != nil {
		return nil, errDetails(err)
	}

	return fn.ParseKubeObject(b)
}
