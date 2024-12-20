package helm

import (
	"fmt"
	"knit/pkg/logging"
	"log"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

func getChart(chartRef *ChartRef, settings *cli.EnvSettings) (*chart.Chart, error) {
	logger := logging.GetInstance()

	cpo := action.ChartPathOptions{
		RepoURL: chartRef.Repository,
		Version: chartRef.Version,
	}

	var chartPath string
	var err error
	locateChart := func() error {
		var err error
		chartPath, err = cpo.LocateChart(chartRef.Name, settings)
		return err
	}

	err = locateChart()
	if err != nil {
		err := pullChart(logger, settings, chartRef)
		if err != nil {
			log.Fatalln(fmt.Errorf("failed to pull chart: %w", err))
			return nil, err
		}

		err = locateChart()
		if err != nil {
			log.Fatalln(fmt.Errorf("failed to locate chart: %w", err))
			return nil, err
		}
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to load chart: %w", err))
		return nil, err
	}

	return chart, nil
}

func RunTemplate(chartRef *ChartRef, values map[string]interface{}, releaseName, namespace string, apiVersions []string) (*release.Release, error) {
	settings := cli.New()
	chart, err := getChart(chartRef, settings)
	if err != nil {
		return nil, err
	}

	registryClient, err := newRegistryClient(settings, false)
	if err != nil {
		return nil, fmt.Errorf("failed to created registry client: %w", err)
	}

	kubeVersion, err := chartutil.ParseKubeVersion("1.31.0")
	if err != nil {
		return nil, err
	}

	client := action.NewInstall(&action.Configuration{})
	client.SetRegistryClient(registryClient)
	client.DryRun = true
	client.ClientOnly = true
	client.ReleaseName = releaseName
	client.Namespace = namespace
	client.KubeVersion = kubeVersion
	if len(apiVersions) == 0 {
		client.APIVersions = chartutil.DefaultVersionSet
	} else {
		client.APIVersions = apiVersions
	}
	client.IncludeCRDs = true
	client.DisableHooks = true

	rel, err := client.Run(chart, values)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to run install: %w", err))
		return nil, err
	}

	if rel == nil {
		return nil, fmt.Errorf("nil release")
	}

	return rel, nil
}
