package main

import (
	"twowls.org/patchwork/commons/extension"
	"twowls.org/patchwork/commons/utils/singleton"
)

type zerologPluginInfo struct {
	loggerExt singleton.Lazy[*zerologExtension]
}

func (p *zerologPluginInfo) Description() string {
	return "zerolog rich console logger"
}

func (p *zerologPluginInfo) DefaultExtension() extension.Extension {
	return p.loggerExt.Instance()
}

func PluginInfo() (extension.PluginInfo, error) {
	return &zerologPluginInfo{
		loggerExt: singleton.NewLazy(func() *zerologExtension {
			return &zerologExtension{}
		}),
	}, nil
}
