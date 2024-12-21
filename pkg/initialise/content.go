package initialise

const helmBaseFileContent = `"""
This file was generated by knit. DO NOT EDIT.
To update the repository or version of this chart, re-run the 'knit add helm' command.
"""
import kcl_plugin.helm
import knit.kubernetes

template = lambda chart: Chart -> [kubernetes.Manifest] {
    helm.template(chart)
}

schema Chart:
    repository: str
    name: str
    version: str
    releaseName: str
    namespace: str = ""
    values: Values
    capabilities?: Capabilities

schema Capabilities:
    apiVersions?: [str]

schema Values:
    [str]: any
`

const kubernetesBaseFileContent = `"""
This file was generated by knit. DO NOT EDIT.
To update the repository or version of this chart, re-run the 'knit add helm' command.
"""
schema Manifest:
    apiVersion: str
    kind: str
    metadata: ManifestMetadata
    [str]: any

schema ManifestMetadata:
    [str]: any
    name?: str
    namespace?: str
    generateName?: str
    labels?: {str: str}
    annotations?: {str: str}
`
