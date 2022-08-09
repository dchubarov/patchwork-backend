package main

import (
	"twowls.org/patchwork/commons/extension"
)

func PluginInfo() (extension.PluginInfo, error) {
	return &zerologPluginInfo{}, nil
}

type zerologPluginInfo struct{}

func (info *zerologPluginInfo) Description() string {
	return "zerolog rich console logger"
}

func (info *zerologPluginInfo) DefaultExtension() extension.Extension {
	return &zerologExtension{}
}
