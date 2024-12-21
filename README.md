# knit

knit is a tool containing the missing pieces for using KCL to define kubernetes manifests when using the [rendered manifests pattern](https://akuity.io/blog/the-rendered-manifests-pattern).

It knits together **Helm** and **KCL** to help keep code manageable and most imortantly, type-safe.

### Helm

knit incorporates **Helm** so that you can leverage existing charts in KCL without having to run a pipeline such as `helm template | kcl import`. Instead, you can add helm charts using `knit add helm` which will create a vendored KCL module for the chart which you can import. Then you can render out the chart using the `helm.template()` function within your KCL when `knit render` is run.

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

You can also set values for the helm chart, for example:
```kcl
import vendored.helm.podinfo
import knit.helm
import manifests

_chart = podinfo.Chart {
    releaseName = "my-app"
    values = podinfo.Values {
        ingress = podinfo.ValuesIngress {
            enabled = True
        }
    }
}

# We use the builtin yaml_stream function here because we cannot put a list and a field at the top-level in the same file
manifests.yaml_stream(helm.template(_chart))

```

> The examples above are just minimal examples. Break things up into modules as required and utilise KCL's language features to make things work for you.

You can view this example [here](example/).

## Notes

### KCL Plugins
knit uses KCL plugins. These plugins only work within KCL VMs where they are explicitly imported. This is done within the `knit render` command which uses the KCL sdk to instantiate a KCL interpreter with the custom plugins. This means that you will encounter errors if using `kcl run` on files which use the custom plugin functions.
