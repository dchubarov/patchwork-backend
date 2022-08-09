package main

import (
	"twowls.org/patchwork/commons/extension"
	"twowls.org/patchwork/commons/utils"
)

func PluginInfo() (extension.PluginInfo, error) {
	return &zerologPluginInfo{}, nil
}

type zerologPluginInfo struct {
	defaultExt utils.Singleton[extension.Extension]
}

func (p *zerologPluginInfo) Description() string {
	return "zerolog rich console logger"
}

func (p *zerologPluginInfo) DefaultExtension() extension.Extension {
	return p.defaultExt.Instance(func() extension.Extension {
		return &zerologExtension{}
	})
}
