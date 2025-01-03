package helm

const chartFileContent = `"""
This file was generated by knit. DO NOT EDIT.
To update the repository or version of this chart, re-run the 'knit add helm' command.
"""
import knit.helm

schema Chart(helm.Chart):
    repository: "%s" = "%s"
    name: "%s" = "%s"
    version: "%s" = "%s"
    values: helm.Values = Values {}
`
