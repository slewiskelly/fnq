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
	_ "embed"
	"os"

	"cuelang.org/go/cue/cuecontext"
	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/slewiskelly/fnq"
)

func main() {
	v := cuecontext.New().CompileString( /* ... */ )

	err := fn.AsMain(fn.Chain(
		fnq.Generate(v),  // Generate additional resources, adding them to the resource list.
		fnq.Transform(v), // Transform resources within the resource list.
		fnq.Validate(v),  // Validate all resources within the resource list.
	))
	if err != nil {
		os.Exit(1)
	}
}
```
