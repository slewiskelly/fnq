# CRD

This example shows how to create a CRD which is used to generate a list of
deployable Kubernetes resources.

## Summary

Files of interest are as follows:

- `./cmd/fnq/main.go`
  - Entrypoint used to first validate and then generate resources based on an
    input `Application` CRD
- `./data/resourcelist.yaml`
  - A Kubernetes resource list containing a single, simple `Application` CRD
- `./mod/application.cue`
  - Defines a `Validator` and `Generator` to validate and generate resources
    from an input `Application` CRD

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
   contents of the executable will output a generated list of resources

```shell
cat ./data/resourcelist.yaml | CUE_REGISTRY=127.0.0.1:5001 ./bin/fnq
```
