package helm

import (
	"os"
	"path/filepath"
)

type ChartFiles struct {
	chart      *Chart
	path       string
	chartPath  string
	isTempPath bool
}

func newChartFiles(chart *Chart, path string, isTempPath bool) (*ChartFiles, error) {
	return &ChartFiles{
		chart:      chart,
		path:       path,
		isTempPath: isTempPath,
		chartPath:  filepath.Join(path, filepath.Clean(chart.chart.Name)),
	}, nil
}

func (c *ChartFiles) ChartPath() string {
	return c.chartPath
}

func (c *ChartFiles) Close() error {
	if c.isTempPath && c.path != "" {
		return os.RemoveAll(c.path)
	}
	return nil
}
