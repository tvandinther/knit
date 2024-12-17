package helm

import (
	"bytes"
	"fmt"
	"knit/pkg/logging"
	"log"
	"sort"
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/releaseutil"
)

func RunPull(logger *log.Logger, settings *cli.EnvSettings, repository, chartRef, version string) error {

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

func RunTemplateNew() error {
	logger := logging.GetInstance()

	registryClient, err := newRegistryClient(&cli.EnvSettings{}, false)
	if err != nil {
		return fmt.Errorf("failed to created registry client: %w", err)
	}

	kubeVersion, err := chartutil.ParseKubeVersion("1.31.0")
	if err != nil {
		return err
	}
	
	client := action.NewInstall(&action.Configuration{})
	client.SetRegistryClient(registryClient)
	client.DryRun = true
	client.ReleaseName = "todo"
	client.Replace = true // TODO: Why?
	client.ClientOnly = true
	client.APIVersions = chartutil.DefaultVersionSet // TODO: What's this?
	client.IncludeCRDs = true
	client.KubeVersion = kubeVersion

	cpo := action.ChartPathOptions{
		RepoURL: "https://stefanprodan.github.io/podinfo",
		Version: "6.7.1",
	}

	var chartPath string
	locateChart := func() {
		chartPath, err = cpo.LocateChart("podinfo", &cli.EnvSettings{})
	}

	locateChart()
	if err != nil {
		err := RunPull(logger, &cli.EnvSettings{}, cpo.RepoURL, "podinfo", cpo.Version)
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

	values, err := chartutil.ReadValuesFile(chartPath + "/values.yaml")
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to load values: ", err))
		return err
	}

	rel, err := client.Run(chart, values)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to run install: ", err))
		return err
	}

	if rel == nil {
		return fmt.Errorf("nil release")
	}

	print(rel)
	
	return nil
}

func print(rel *release.Release) {
	if rel != nil {
		var manifests bytes.Buffer
		fmt.Fprintln(&manifests, strings.TrimSpace(rel.Manifest))
		// if !client.DisableHooks {
		// 	fileWritten := make(map[string]bool)
		// 	for _, m := range rel.Hooks {
		// 		if skipTests && isTestHook(m) {
		// 			continue
		// 		}
		// 		if client.OutputDir == "" {
		// 			fmt.Fprintf(&manifests, "---\n# Source: %s\n%s\n", m.Path, m.Manifest)
		// 		} else {
		// 			newDir := client.OutputDir
		// 			if client.UseReleaseName {
		// 				newDir = filepath.Join(client.OutputDir, client.ReleaseName)
		// 			}
		// 			_, err := os.Stat(filepath.Join(newDir, m.Path))
		// 			if err == nil {
		// 				fileWritten[m.Path] = true
		// 			}

		// 			err = writeToFile(newDir, m.Path, m.Manifest, fileWritten[m.Path])
		// 			if err != nil {
		// 				return err
		// 			}
		// 		}

		// 	}
		// }

		// if we have a list of files to render, then check that each of the
		// provided files exists in the chart.
		// if len(showFiles) > 0 {
		if false {
			// This is necessary to ensure consistent manifest ordering when using --show-only
			// with globs or directory names.
			splitManifests := releaseutil.SplitManifests(manifests.String())
			manifestsKeys := make([]string, 0, len(splitManifests))
			for k := range splitManifests {
				manifestsKeys = append(manifestsKeys, k)
			}
			sort.Sort(releaseutil.BySplitManifestsOrder(manifestsKeys))

			// manifestNameRegex := regexp.MustCompile("# Source: [^/]+/(.+)")
			var manifestsToRender []string
			// for _, f := range showFiles {
			// 	missing := true
			// 	// Use linux-style filepath separators to unify user's input path
			// 	f = filepath.ToSlash(f)
			// 	for _, manifestKey := range manifestsKeys {
			// 		manifest := splitManifests[manifestKey]
			// 		submatch := manifestNameRegex.FindStringSubmatch(manifest)
			// 		if len(submatch) == 0 {
			// 			continue
			// 		}
			// 		manifestName := submatch[1]
			// 		// manifest.Name is rendered using linux-style filepath separators on Windows as
			// 		// well as macOS/linux.
			// 		manifestPathSplit := strings.Split(manifestName, "/")
			// 		// manifest.Path is connected using linux-style filepath separators on Windows as
			// 		// well as macOS/linux
			// 		manifestPath := strings.Join(manifestPathSplit, "/")

			// 		// if the filepath provided matches a manifest path in the
			// 		// chart, render that manifest
			// 		if matched, _ := filepath.Match(f, manifestPath); !matched {
			// 			continue
			// 		}
			// 		manifestsToRender = append(manifestsToRender, manifest)
			// 		missing = false
			// 	}
			// 	if missing {
			// 		return fmt.Errorf("could not find template %s in chart", f)
			// 	}
			// }
			for _, m := range manifestsToRender {
				out := log.Writer()
				fmt.Fprintf(out, "---\n%s\n", m)
			}
		} else {
			out := log.Writer()
			fmt.Fprintf(out, "%s", manifests.String())
		}
	}
}
