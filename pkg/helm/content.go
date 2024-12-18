package helm

const chartFileContent = `schema Chart:
    repository: "%s" = "%s"
    name: "%s" = "%s"
    version: "%s" = "%s"
    values: Values = Values{}
`
