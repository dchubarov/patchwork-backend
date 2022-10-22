package main

import (
	"twowls.org/patchwork/commons/extension"
	"twowls.org/patchwork/plugin/database/mongodb/client"
)

type mongodbPlugin struct {
	clientExtension *client.ClientExtension
}

func (p *mongodbPlugin) Description() string {
	return "Database operations backed by MongoDB"
}

func (p *mongodbPlugin) DefaultExtension() extension.Extension {
	return p.clientExtension
}

func PluginInfo() (extension.PluginInfo, error) {
	return &mongodbPlugin{new(client.ClientExtension)}, nil
}
