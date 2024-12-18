package helm

import (
	"fmt"
	"log"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)


type ChartRef struct {
	Repository string
	Name string
	Version string
}

func pullChart(logger *log.Logger, settings *cli.EnvSettings, chartRef *ChartRef) error {

	actionConfig, err := initActionConfig(settings, logger)
	if err != nil {
		return fmt.Errorf("failed to init action config: %w", err)
	}

	registryClient, err := newRegistryClient(settings, false)
	if err != nil {
		return fmt.Errorf("failed to created registry client: %w", err)
	}
	actionConfig.RegistryClient = registryClient

	logger.Println(settings.RepositoryCache)
	pullClient := action.NewPullWithOpts(action.WithConfig(actionConfig))
	pullClient.RepoURL = chartRef.Repository
	pullClient.DestDir = settings.RepositoryCache
	pullClient.Settings = settings
	pullClient.Version = chartRef.Version

	result, err := pullClient.Run(chartRef.Name)
	if err != nil {
		return fmt.Errorf("failed to pull chart: %w", err)
	}

	logger.Printf("%+v", result)

	return nil
}
