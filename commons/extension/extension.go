package extension

// Extension provides methods to integrate with application
type Extension interface {
	// Configure called when application needs to integrate the extension
	Configure(options *Options) error
}

// PluginInfo provides plug-in information
type PluginInfo interface {
	// Description returns plug-in's description
	Description() string
	// DefaultExtension returns default extension provided by plug-in, or nil if there is none
	DefaultExtension() Extension
}

// PluginEntrypoint name of entrypoint function
const PluginEntrypoint = "PluginInfo"
