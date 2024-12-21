# knit

knit is a tool containing the missing pieces for using KCL to define kubernetes manifests when using the [rendered manifests pattern](https://akuity.io/blog/the-rendered-manifests-pattern).

It knits together **Helm** and **Kustomize**, with **KCL** to help keep code manageable and most imortantly, type-safe.

### KCL

[KCL](https://www.kcl-lang.io/) is a configuration language and a CNCF sandbox project. It is used as the way to define your data structures empowered with modern language features to enable you to write DRY code.

### Helm

knit incorporates **Helm** so that you can leverage existing charts in KCL without having to run a pipeline such as `helm template | kcl import`. Instead, you can add helm charts using `knit add helm` which will create a vendored KCL module for the chart which you can import. Then you can render out the chart using the `helm.template()` function within your KCL when `knit render` is run.

### Kustomize

knit also includes functions to interact with **Kustomize** to enable you to use existing transformations and overlays without requiring a full migration to KCL templating or to empower your KCL code with Kustomize features. Use the `kustomize.build()` function to run a kustomization.

## Quick Start

Start by initialising a project in a directory. This command creates KCL module files just like `kcl mod init` does.
```sh
knit init
```

Add a helm chart to the project.
```sh
knit add helm https://stefanprodan.github.io/podinfo podinfo --version 6.7.1
```

In `main.k`, try rendering the default chart.
```kcl
# main.k
import vendored.helm.podinfo
import knit.helm

[manifest for manifest in helm.template(podinfo.Chart{})]

```

```sh
knit render
```

You can also set values for the Helm chart, for example:
```kcl
# main.k
import vendored.helm.podinfo
import knit.helm
import knit.kustomize
import manifests

_chart = podinfo.Chart {
    releaseName = "my-app"
    values = podinfo.Values {
        ingress = podinfo.ValuesIngress {
            enabled = True
        }
    }
}

_manifests = helm.template(_chart)
```

You can post-process the helm chart with Kustomize. The kustomization provided does not include resource relative path names. Instead, the resources are given as a list of KCL maps in an additional argument. For example, to set the namespace of all resources and to change the tag of the image you can do the following:
```kcl
# main.k continued...
_kustomization = kustomize.Kustomization {
    namespace = "team-a"
    images = [{
        name = "ghcr.io/stefanprodan/podinfo"
        digest = "sha256:862ca45e61b32392f7941a1bdfdbe5ff8b6899070135f1bdca1c287d0057fc94"
    }]
}

_kustomized = kustomize.build(_kustomization, _manifests)
```

Then collect the list of manifests into a stream of YAML documents.
```kcl
# main.k continued...
# We use the builtin yaml_stream function here because we cannot put a list and a field at the top-level in the same file
manifests.yaml_stream(_kustomized)

```

Finally, run render again and inspect the changes kustomize has made:
```sh
knit render | grep -e "kind:" -e "image:" -e "namespace:"
```

> The examples above are just minimal examples. Break things up into modules as required and utilise KCL's language features to make things work for you.

You can view this example [here](example/).

## Reference

### Functions
The table below shows the KCL functions available to use from the `knit` module.

| Import | Function | Description |
| --- | --- | --- |
| `import knit.helm` | `helm.template(chart: helm.Chart) -> [kubernetes.Manifest]` | Runs an equivalent of `helm template` on the given chart and returns a list of kubernetes manifests. |
| `import knit.kustomize` | `kustomize.build(kustomization: kustomize.Kustomization, resources: [{str: any}]) -> [kubernetes.Manifest]` | Runs an equivalent of `kustomize build` with the given base kustomization and resources. |

### Schemas
The table below shows the KCL schemas available to use from the `knit` module.

| Import | Schema | Description |
| --- | --- | --- |
| `import knit.helm` | `helm.Chart` | Represents a helm chart to be parsed by the `helm.template` function. |
| `import knit.kustomize` | `kustomize.Kustomization` | The schema for a valid `kustomization.yaml`. |
| `import knit.kubernetes` | `kubernetes.Manifest` | A minimal schema for a Kubernetes resource manifest. |

## Notes

### Kustomize
Mutations offered by Kustomize can often be performed using simple KCL. For example, using the namespace transformer can be done simply using KCL dict unions. The following shows setting the namespace on a list of `Manifests` using both approaches.
```kcl
# Using Kustomize
_kustomized = kustomize.build({namespace: "team-a"}, _manifests)

# Using KCL dict unions
_kustomized = [m {metadata.namespace = "team-a"} for m in _manifests]
```
More complex transformation may be better off left to Kustomize to perform. You have the freedom to combine the power of both approaches.

### KCL Plugins
knit uses KCL plugins. These plugins only work within KCL VMs where they are explicitly imported. This is done within the `knit render` command which uses the KCL sdk to instantiate a KCL interpreter with the custom plugins. This means that you will encounter errors if using `kcl run` on files which use the custom plugin functions.
