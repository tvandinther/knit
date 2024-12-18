# knit

knit is a tool containing the missing pieces for using KCL to define kubernetes manifests when using the [rendered manifests pattern](https://akuity.io/blog/the-rendered-manifests-pattern).

It knits together **Helm**, **git** and **KCL** to help keep code manageable and most imortantly, type-safe.

### Helm

knit incorporates **Helm** so that you can leverage existing charts in KCL without having to run a pipeline such as `helm template | kcl import`. Instead, you can add helm charts using `knit add helm` which will create a vendored KCL module for the chart which you can import. Then you can render out the chart using the `helm.Template()` function within your KCL when `knit render` is run.

### git (in-progress)

knit also incorporates **git** to take the hassle out of automating git actions related to the rendered manifests pattern. It does this by offering the `knit ___` command to run a multi-step process of branching, rendering, staging, and committing the new manifests.

## Examples

```sh
knit add helm https://stefanprodan.github.io/podinfo podinfo --version 6.7.1
```
