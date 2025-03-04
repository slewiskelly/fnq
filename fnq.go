// Package fnq accommodates the creation of KRM functions powered by [CUE].
//
// See [usage] and [examples] for more information on how to create a function
// and [KRM Functions Specification] for more information with regards to the
// specification.
//
// [CUE]: https://cuelang.org/
// [KRM Functions Specification]: https://github.com/kubernetes-sigs/kustomize/blob/master/cmd/config/docs/api-conventions/functions-spec.md
// [examples]: https://github.com/slewiskelly/fnq/examples
// [usage]: https://github.com/slewiskelly/fnq/docs/references/usage
package fnq

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
)

type processor struct {
	v cue.Value
}

func errDetails(e error) error {
	qe, ok := e.(errors.Error)
	if !ok {
		return e
	}

	var msgs []string

	qe = errors.Sanitize(qe)

	for _, err := range errors.Errors(qe) {
		f, a := err.Msg()
		msgs = append(msgs, fmt.Sprintf("%s: %s", strings.Join(err.Path()[2:], "."), fmt.Sprintf(f, a...)))
	}

	return errors.New(strings.Join(msgs, "\n"))
}

func gvk(obj *fn.KubeObject) string {
	gvk := obj.GroupVersionKind()

	return fmt.Sprintf("%s_%s_%s", gvk.Group, gvk.Version, gvk.Kind)
}

var errNotFound = errors.New("not found")
