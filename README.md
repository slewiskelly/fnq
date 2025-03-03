# fnq

Package `fnq` accommodates the creation of KRM functions powered by [CUE](https://cuelang.org/).

> [!WARNING]
> ___This is a work in progress and is considered experimental.___

## Usage

```go
import "github.com/slewiskelly/fnq"
```

```go
package main

import (
	"context"
	_ "embed"
	"os"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/slewiskelly/fnq"
)

func main() {
	ctx := context.Background()

	err := fn.AsMain(fn.Chain(
		fnq.Generate(ctx, module), // Generate additional resources, adding them to the resource list.
		fnq.Transform(ctx, modle), // Transform resources within the resource list.
		fnq.Validate(ctx, module), // Validate all resources within the resource list.
	))
	if err != nil {
		os.Exit(1)
	}
}

var module string // Set via linker flags.
```
