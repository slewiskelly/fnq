package fnq

import (
	"context"
	"fmt"
	"slices"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/encoding/yaml"
	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"

	"github.com/slewiskelly/fnq/internal/pkg/mod"
)

// Generate returns a ResourceListProcessorFunc for generating resources
// within a resource list.
//
// Resources will be generated from other resources within the resource list
// according to a corresponding resource generator in the given module.
//
// See the [CUE reference] page for how validators should be expressed.
//
// A panic will occur if there is an error while retrieving the module.
//
// [CUE reference]: https://github.com/slewiskelly/fnq/docs/references/cue.md
func Generate(ctx context.Context, module string) fn.ResourceListProcessorFunc {
	v, err := mod.Get(ctx, module)
	if err != nil {
		panic(err)
	}

	return (&processor{v: v}).generate
}

func (p *processor) generate(rl *fn.ResourceList) (bool, error) {
	var items fn.KubeObjects
	var results fn.Results

	for _, obj := range rl.Items.WhereNot(fn.HasAnnotations(map[string]string{"no-generate": "true"})) {
		i, err := p.gen(obj)
		if err != nil {
			if errors.Is(err, errNotFound) {
				continue
			}

			results = append(results, fn.ErrorConfigObjectResult(err, obj))
			continue
		}

		rl.Results.Infof("Generated %d resources from %s", len(i), obj.GetId())

		items = append(items, i...)
	}

	rl.Items = slices.Concat(rl.Items.WhereNot(fn.HasAnnotations(map[string]string{"local-only": "true"})), items)
	rl.Results = append(rl.Results, results...)

	if len(results) > 0 {
		return false, errors.New("generation failed")
	}

	return true, nil
}

func (p *processor) gen(obj *fn.KubeObject) (fn.KubeObjects, error) {
	grp, ver, kind := gvk(obj)

	v := p.v.LookupPath(cue.ParsePath(fmt.Sprintf("Generators[%q][%q][%q]", grp, ver, kind)))
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
	if err := v.Validate(cue.Final(), cue.Hidden(true)); err != nil {
		return nil, errDetails(err)
	}

	iter, err := v.Fields(cue.Hidden(true))
	if err != nil {
		return nil, errDetails(err)
	}

	for iter.Next() {
		if iter.Selector().String() != "_$$resources" {
			continue
		}

		rs, err := iter.Value().List()
		if err != nil {
			return nil, errDetails(err)
		}

		b, err := yaml.EncodeStream(rs)
		if err != nil {
			return nil, err
		}

		return fn.ParseKubeObjects(b)
	}

	return nil, nil
}
