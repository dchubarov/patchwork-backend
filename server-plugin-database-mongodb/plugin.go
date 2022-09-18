package main

import (
	"twowls.org/patchwork/commons/extension"
)

type mongodbPlugin struct {
	clientExtension *ClientExtension
}

func (p *mongodbPlugin) Description() string {
	return "Database operations backed by MongoDB"
}

func (p *mongodbPlugin) DefaultExtension() extension.Extension {
	return p.clientExtension
}

func PluginInfo() (extension.PluginInfo, error) {
	return &mongodbPlugin{new(ClientExtension)}, nil
}
