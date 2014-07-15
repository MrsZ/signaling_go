package newrelic

const (
	PollInterval = 10

	//DefaultAgentGuid is plugin ID in NewRelic.
	//You should not change it unless you want to create your own plugin.
	DefaultAgentGuid = "com.github.msoedov.SignalingRealtime"

	//CurrentAgentVersion is plugin version
	CurrentAgentVersion = "0.0.6"

	//DefaultAgentName in NewRelic GUI. You can change it.
	DefaultAgentName = "Go Realtime app"
)
