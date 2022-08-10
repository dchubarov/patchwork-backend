package main

import (
	"twowls.org/patchwork/commons/extension"
	"twowls.org/patchwork/commons/singleton"
)

type zerologPluginInfo struct {
	loggerExt singleton.S[*zerologExtension]
}

func (p *zerologPluginInfo) Description() string {
	return "zerolog rich console logger"
}

func (p *zerologPluginInfo) DefaultExtension() extension.Extension {
	return p.loggerExt.Instance()
}

func PluginInfo() (extension.PluginInfo, error) {
	return &zerologPluginInfo{
		loggerExt: singleton.Lazy(func() *zerologExtension {
			return &zerologExtension{}
		}),
	}, nil
}
