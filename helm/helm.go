package helm

import (
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
)

var allGetters = getter.All(&cli.EnvSettings{})
