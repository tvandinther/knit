import vendored.helm.podinfo
import kcl_plugin.helm
import manifests

_chart = podinfo.Chart {
    values = podinfo.Values {
        ingress = podinfo.ValuesIngress {
            enabled = True
        }
    }
}

# We use the builtin yaml_stream function here because we cannot put a list and a field at the top-level in the same file
manifests.yaml_stream(helm.template(_chart))
