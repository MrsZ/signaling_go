package newrelic

import (
	"fmt"
	metrics "github.com/yvasiyarov/go-metrics"
	"github.com/yvasiyarov/newrelic_platform_go"
)

type Observable interface {
	GetStats() (metrics.Counter, metrics.Counter)
}

func InitNewrelicAgent(license string, appname string, verbose bool, obs Observable) error {

	if license == "" {
		return fmt.Errorf("Please specify NewRelic license")
	}

	plugin := newrelic_platform_go.NewNewrelicPlugin(CurrentAgentVersion, license, PollInterval)
	component := newrelic_platform_go.NewPluginComponent(DefaultAgentName, DefaultAgentGuid)
	plugin.AddComponent(component)

	m := &MembersMetrica{obs}

	component.AddMetrica(m)

	plugin.Verbose = verbose
	plugin.Run()
	return nil
}
