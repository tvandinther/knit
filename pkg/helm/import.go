package helm

func Import(chartref *ChartRef, directory string) error {
	err := getValues(chartref, directory)

	return err
}
