package main

import (
	"context"
	"os"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/slewiskelly/fnq"
)

func main() {
	ctx := context.Background()

	err := fn.AsMain(fnq.Transform(ctx, module))
	if err != nil {
		os.Exit(1)
	}
}

var module string // Set via linker flags.
