package helm

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"github.com/pkg/errors"
// 	"helm.sh/helm/pkg/downloader"
// 	"helm.sh/helm/pkg/getter"
// 	"helm.sh/helm/v3/pkg/action"
// 	"helm.sh/helm/v3/pkg/chart/loader"
// 	"helm.sh/helm/v3/pkg/cli/values"
// 	"helm.sh/helm/v3/pkg/release"
// )

// func runInstall(args []string, client *action.Install, valueOpts *values.Options, out io.Writer) (*release.Release, error) {
// 	name, chart, err := client.NameAndChart(args)
// 	if err != nil {
// 		return nil, err
// 	}
// 	client.ReleaseName = name

// 	cp, err := client.ChartPathOptions.LocateChart(chart, settings)
// 	if err != nil {
// 		return nil, err
// 	}

// 	p := getter.All(settings)
// 	vals, err := valueOpts.MergeValues(p)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Check chart dependencies to make sure all are present in /charts
// 	chartRequested, err := loader.Load(cp)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := checkIfInstallable(chartRequested); err != nil {
// 		return nil, err
// 	}

// 	if chartRequested.Metadata.Deprecated {
// 		warning("This chart is deprecated")
// 	}

// 	if req := chartRequested.Metadata.Dependencies; req != nil {
// 		// If CheckDependencies returns an error, we have unfulfilled dependencies.
// 		// As of Helm 2.4.0, this is treated as a stopping condition:
// 		// https://github.com/helm/helm/issues/2209
// 		if err := action.CheckDependencies(chartRequested, req); err != nil {
// 			err = errors.Wrap(err, "An error occurred while checking for chart dependencies. You may need to run `helm dependency build` to fetch missing dependencies")
// 			if client.DependencyUpdate {
// 				man := &downloader.Manager{
// 					Out:              out,
// 					ChartPath:        cp,
// 					Keyring:          client.ChartPathOptions.Keyring,
// 					SkipUpdate:       false,
// 					Getters:          p,
// 					RepositoryConfig: settings.RepositoryConfig,
// 					RepositoryCache:  settings.RepositoryCache,
// 					Debug:            settings.Debug,
// 					RegistryClient:   client.GetRegistryClient(),
// 				}
// 				if err := man.Update(); err != nil {
// 					return nil, err
// 				}
// 				// Reload the chart with the updated Chart.lock file.
// 				if chartRequested, err = loader.Load(cp); err != nil {
// 					return nil, errors.Wrap(err, "failed reloading chart after repo update")
// 				}
// 			} else {
// 				return nil, err
// 			}
// 		}
// 	}

// 	client.Namespace = settings.Namespace()

// 	// Validate DryRunOption member is one of the allowed values
// 	if err := validateDryRunOptionFlag(client.DryRunOption); err != nil {
// 		return nil, err
// 	}

// 	// Create context and prepare the handle of SIGTERM
// 	ctx := context.Background()
// 	ctx, cancel := context.WithCancel(ctx)

// 	// Set up channel on which to send signal notifications.
// 	// We must use a buffered channel or risk missing the signal
// 	// if we're not ready to receive when the signal is sent.
// 	cSignal := make(chan os.Signal, 2)
// 	signal.Notify(cSignal, os.Interrupt, syscall.SIGTERM)
// 	go func() {
// 		<-cSignal
// 		fmt.Fprintf(out, "Release %s has been cancelled.\n", args[0])
// 		cancel()
// 	}()

// 	return client.RunWithContext(ctx, chartRequested, vals)
// }