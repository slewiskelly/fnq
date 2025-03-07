# Transformer

This example shows how to transform a `Deployment` resource.

## Summary

Files of interest are as follows:

- `./cmd/fnq/main.go`
  - Entrypoint used to transform resources within a resource list
- `./data/resourcelist.yaml`
  - A Kubernetes resource list containing a `Deployment` and `Service`
- `./mod/transformer.cue`
  - Defines a `Transformer` that transforms `Deployment` resources

## Usage

1. In a separate terminal window, start a local CUE module registry

```shell
cue mod registry 127.0.0.1:5001
````

2. Publish the module under `./mod` to the local registry

```shell
CUE_REGISTRY=127.0.0.1:5001 VERSION=v0.0.1 make -C ./mod publish
```

3. Build an executable, specifying the version of the module to be used by all
   executions

```shell
VERSION=v0.0.1 make build
```

4. Specifying the local registry from which to copy the module, piping the
   contents of the executable will output a transformed list of resources

```shell
cat ./data/resourcelist.yaml | CUE_REGISTRY=127.0.0.1:5001 ./bin/fnq
```
