package main

import (
	"twowls.org/patchwork/commons/extension"
	"twowls.org/patchwork/plugin/database/mongodb/mongodb"
)

type mongodbPlugin struct {
	clientExtension *mongodb.ClientExtension
}

func (p *mongodbPlugin) Description() string {
	return "mongodb"
}

func (p *mongodbPlugin) DefaultExtension() extension.Extension {
	return p.clientExtension
}

func PluginInfo() (extension.PluginInfo, error) {
	return &mongodbPlugin{new(mongodb.ClientExtension)}, nil
}
