package plugins

import (
	"errors"
	"fmt"
	"os"
	"path"
	"plugin"
	"strings"
	"twowls.org/patchwork/commons/extension"
	"twowls.org/patchwork/server/bootstrap/config"
)

func Load(name string) (extension.PluginInfo, error) {
	location, err := resolve(name)
	if err != nil {
		return nil, err
	}

	module, err := plugin.Open(location)
	if err != nil {
		return nil, err
	}

	entrypoint, err := module.Lookup(extension.PluginEntrypoint)
	if err != nil {
		return nil, err
	}

	if entrypointFunc, ok := entrypoint.(func() (extension.PluginInfo, error)); !ok {
		return nil, errors.New(fmt.Sprintf("invalid plugin entrypoint type: %t", entrypoint))
	} else {
		return entrypointFunc()
	}
}

func resolve(name string) (string, error) {
	if name == "" {
		return "", errors.New("plugin name cannot be empty")
	}

	location := path.Join(config.Values().PluginsDir,
		fmt.Sprintf("plugin_%s_build",
			strings.Replace(name, "-", "_", -1)))

	if stat, err := os.Stat(location); err == nil {
		if (stat.Mode() & os.ModeType) != 0 {
			return "", errors.New("plugin path " + location + " is not a regular file")
		}
	}

	return location, nil
}
