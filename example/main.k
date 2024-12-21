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

_kustomization = kustomize.Kustomization {
    namespace = "team-a"
    images = [{
        name = "ghcr.io/stefanprodan/podinfo"
        digest = "sha256:862ca45e61b32392f7941a1bdfdbe5ff8b6899070135f1bdca1c287d0057fc94"
    }]
}

_kustomized = kustomize.build(_kustomization, _manifests)

manifests.yaml_stream(_kustomized)
