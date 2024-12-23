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

_podinfo = helm.template(_chart)

_podinfoKustomization = kustomize.Kustomization {
    namespace = "team-a"
    images = [{
        name = "ghcr.io/stefanprodan/podinfo"
        digest = "sha256:862ca45e61b32392f7941a1bdfdbe5ff8b6899070135f1bdca1c287d0057fc94"
    }]
}

_kustomizedPodinfo = kustomize.build(_podinfoKustomization, _podinfo)

_argocd = kustomize.build(kustomize.Kustomization {
    resources = ["https://raw.githubusercontent.com/argoproj/argo-cd/v2.7.2/manifests/install.yaml"]
}, [])

# Outputs all kustomized podinfo and Argo CD resource manifests
manifests.yaml_stream(_kustomizedPodinfo + _argocd)
