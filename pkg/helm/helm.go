package helm

import (
	"fmt"
	"knit/pkg/logging"
	"log"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/engine"
	"k8s.io/client-go/rest"
)

func runPull(logger *log.Logger, settings *cli.EnvSettings, repository, chartRef, version string) error {

	actionConfig, err := initActionConfig(settings, logger)
	if err != nil {
		return fmt.Errorf("failed to init action config: %w", err)
	}

	registryClient, err := newRegistryClient(settings, false)
	if err != nil {
		return fmt.Errorf("failed to created registry client: %w", err)
	}
	actionConfig.RegistryClient = registryClient

	pullClient := action.NewPullWithOpts(
		action.WithConfig(actionConfig))
	// client.RepoURL = ""
	pullClient.RepoURL = repository
	pullClient.DestDir = "./"
	pullClient.Settings = settings
	pullClient.Version = version
	pullClient.Untar = true
	pullClient.UntarDir = "./"

	result, err := pullClient.Run(chartRef)
	if err != nil {
		return fmt.Errorf("failed to pull chart: %w", err)
	}

	logger.Printf("%+v", result)

	return nil
}

func RunTemplate() error {
	logger := logging.GetInstance()

	cpo := action.ChartPathOptions{
		RepoURL: "https://stefanprodan.github.io/podinfo",
		Version: "6.7.1",
	}

	var chartPath string
	var err error
	locateChart := func() {
		chartPath, err = cpo.LocateChart("podinfo", &cli.EnvSettings{})
	}

	locateChart()
	if err != nil {
		err := runPull(logger, &cli.EnvSettings{}, cpo.RepoURL, "podinfo", cpo.Version)
		if err != nil {
			log.Fatalln(fmt.Errorf("failed to pull chart: ", err))
			return err
		}	

		locateChart()
		if err != nil {
			log.Fatalln(fmt.Errorf("failed to locate chart: ", err))
			return err
		}
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to load chart: ", err))
		return err
	}

	fileValues, err := chartutil.ReadValuesFile(chartPath + "/values.yaml")
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to load values: ", err))
		return err
	}
	coalescedValues := chartutil.CoalesceTables(chart.Values, fileValues)
	values := map[string]interface{}{
		"Values": coalescedValues,
		"Release": map[string]interface{}{
			"Name": "test",
		},
	}
	eng := engine.New(&rest.Config{})
	eng.Strict = false
	rendered, err := eng.Render(chart, values)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	fmt.Println(rendered)

	return nil
}
