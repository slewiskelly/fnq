package main

import (
	"context"
	"os"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/slewiskelly/fnq"
)

func main() {
	ctx := context.Background()

	err := fn.AsMain(fn.Chain(
		// Validating the CRD before attempting to generate resources from it will
		// provide clearer errors if misconfigured.
		fnq.Validate(ctx, module), fnq.Generate(ctx, module),
	))
	if err != nil {
		os.Exit(1)
	}
}

var module string // Set via linker flags.
