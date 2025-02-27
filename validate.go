package fnq

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/encoding/yaml"
	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
)

// Validate returns a ResourceListProcessorFunc for validating resources
// within a resource list.
//
// Each resource within the resource list will be validated according to a
// corresponding resource validator.
//
// See the [CUE reference] page for how validators should be expressed.
//
// [CUE reference]: https://github.com/slewiskelly/fnq/docs/references/cue.md
func Validate(v cue.Value) fn.ResourceListProcessorFunc {
	if err := v.Err(); err != nil {
		panic(fmt.Errorf("invalid validator: %w", err))
	}

	return (&processor{v: v}).validate
}

func (p *processor) validate(rl *fn.ResourceList) (bool, error) {
	var results fn.Results

	for _, obj := range rl.Items.WhereNot(fn.HasAnnotations(map[string]string{"no-validate": "true"})) {
		if err := p.vet(obj); err != nil {
			if errors.Is(err, errNotFound) {
				continue
			}

			results = append(results, fn.ErrorConfigObjectResult(err, obj))
			continue
		}
	}

	rl.Results = append(rl.Results, results...)

	if len(results) > 0 {
		return false, errors.New("validation(s) failed")
	}

	return true, nil
}

func (p *processor) vet(obj *fn.KubeObject) error {
	v := p.v.LookupPath(cue.ParsePath("Validators")).LookupPath(cue.ParsePath(gvk(obj)))
	if err := v.Err(); err != nil {
		if !v.Exists() {
			return errNotFound
		}

		return errDetails(err)
	}

	a, err := yaml.Extract("", obj.String())
	if err != nil {
		return errDetails(err)
	}

	w := v.Context().BuildFile(a)
	if err := w.Err(); err != nil {
		return errDetails(err)
	}

	if err := v.Unify(w).Validate(cue.Hidden(false)); err != nil {
		return errDetails(err)
	}

	return nil
}
