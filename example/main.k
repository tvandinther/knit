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
_kustomized = kustomize.build({namespace: "team-a", images: [{name: "ghcr.io/stefanprodan/podinfo", newTag: "modified"}]}, _manifests)

# We use the builtin yaml_stream function here because we cannot put a list and a field at the top-level in the same file
manifests.yaml_stream(_kustomized)
