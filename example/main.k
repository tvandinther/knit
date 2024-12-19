import vendored.helm.podinfo
import kcl_plugin.helm

[manifest for manifest in helm.template(podinfo.Chart{})]
